package middleware

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var RefreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))

const TokenExpireDuration = time.Hour * 2
const RefreshTokenExpireDuration = time.Hour * 24 * 7 // Refresh token valid for 7 days

// UserClaims represents user JWT token claims
type UserClaims struct {
	ID      uint64   `json:"id"`
	Account string   `json:"account"`
	Roles   []string `json:"roles"`
	jwt.RegisteredClaims
}

// VoterClaims represents voter JWT token claims
type VoterClaims struct {
	ID      uint64    `json:"id"`
	SessionID  uuid.UUID `json:"sessionID"`
	IsVoted bool      `json:"isVoted"`
	jwt.RegisteredClaims
}

// GenUserToken generates user JWT token
func GenUserToken(Id uint64, account string, roles []string) (string, string, error) {
	// Access Token
	accessClaims := UserClaims{
		ID:      Id,
		Account: account,
		Roles:   roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}

	return GenToken(accessClaims)
}

// GenVoterToken generates voter JWT token
func GenVoterToken(Id uint64, sessionID uuid.UUID, isVoted bool) (string, string, error) {
	accessClaims := VoterClaims{
			ID:      Id,
			SessionID:  sessionID,
			IsVoted: isVoted,
			RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}

	return GenToken(accessClaims)
}

// GenToken generates JWT token
func GenToken(accessClaims jwt.Claims) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(SecretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpireDuration)),
		Issuer:    os.Getenv("APP_NAME"),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(RefreshSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ParseUserToken parses user JWT token
func ParseUserToken(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := parseToken(tokenString, claims)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseVoterToken parses voter JWT token
func ParseVoterToken(tokenString string) (*VoterClaims, error) {
	claims := &VoterClaims{}
	token, err := parseToken(tokenString, claims)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*VoterClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// parseToken is a generic function to parse JWT token
func parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
}

// ParseRefreshToken parses and validates refresh token
func ParseRefreshToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return RefreshSecretKey, nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return errors.New("refresh token expired")
		}
		return nil
	}

	return errors.New("invalid refresh token")
}

// JWTAuthMiddleware is an optional JWT middleware (for GraphQL)
// It doesn't block requests, but attempts to parse and set all available token claims
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Attempt to get token from Header
		var headerToken string
		if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
			if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
				headerToken = parts[1]
			}
		}

		// Attempt to parse user-token (from Header or Cookie)
		userTokenString := headerToken
		if userTokenString == "" {
			if userToken, err := c.Cookie("user-token"); err == nil && userToken != "" {
				userTokenString = userToken
			}
		}
		if userTokenString != "" {
			if userClaims, err := ParseUserToken(userTokenString); err == nil {
				c.Set("userId", userClaims.ID)
				c.Set("userAccount", userClaims.Account)
				c.Set("userRoles", userClaims.Roles)
				c.Set("hasUserToken", true)
			}
		}

		// Attempt to parse voter-token (from Cookie, since Header can only have one)
		if voterToken, err := c.Cookie("voter-token"); err == nil && voterToken != "" {
			if voterClaims, err := ParseVoterToken(voterToken); err == nil {
				c.Set("voterId", voterClaims.ID)
				c.Set("voterVoteId", voterClaims.SessionID)
				c.Set("voterIsVoted", voterClaims.IsVoted)
				c.Set("hasVoterToken", true)
			}
		}

		c.Next()
	}
}

// RequireUserToken validates user token in GraphQL resolver
// Returns UserClaims or error
func RequireUserToken(c *gin.Context) (*UserClaims, error) {
	// First check if already set by middleware
	if hasToken, exists := c.Get("hasUserToken"); exists && hasToken.(bool) {
		return &UserClaims{
			ID:      c.GetUint64("userId"),
			Account: c.GetString("userAccount"),
			Roles:   c.GetStringSlice("userRoles"),
		}, nil
	}

	// Otherwise attempt to parse manually
	var tokenString string
	if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}
	if tokenString == "" {
		if userToken, err := c.Cookie("user-token"); err == nil && userToken != "" {
			tokenString = userToken
		}
	}

	if tokenString == "" {
		return nil, errors.New("user token not found")
	}

	return ParseUserToken(tokenString)
}

// RequireVoterToken validates voter token in GraphQL resolver
// Returns VoterClaims or error
func RequireVoterToken(c *gin.Context) (*VoterClaims, error) {
	// First check if already set by middleware
	if hasToken, exists := c.Get("hasVoterToken"); exists && hasToken.(bool) {
		sessionID, _ := c.Get("voterVoteId")
		return &VoterClaims{
			ID:      c.GetUint64("voterId"),
			SessionID:  sessionID.(uuid.UUID),
			IsVoted: c.GetBool("voterIsVoted"),
		}, nil
	}

	// Otherwise attempt to parse manually
	var tokenString string
	if voterToken, err := c.Cookie("voter-token"); err == nil && voterToken != "" {
		tokenString = voterToken
	}

	if tokenString == "" {
		return nil, errors.New("voter token not found")
	}

	return ParseVoterToken(tokenString)
}

// GetOptionalUserToken attempts to get user token in GraphQL resolver (not required)
// Returns UserClaims or nil
func GetOptionalUserToken(c *gin.Context) *UserClaims {
	claims, _ := RequireUserToken(c)
	return claims
}

// GetOptionalVoterToken attempts to get voter token in GraphQL resolver (not required)
// Returns VoterClaims or nil
func GetOptionalVoterToken(c *gin.Context) *VoterClaims {
	claims, _ := RequireVoterToken(c)
	return claims
}
