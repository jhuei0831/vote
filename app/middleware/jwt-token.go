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
const RefreshTokenExpireDuration = time.Hour * 24 * 7 // Refresh Token 有效期 7 天

// UserClaims 用戶 JWT 令牌
type UserClaims struct {
	ID      uint64   `json:"id"`
	Account string   `json:"account"`
	Roles   []string `json:"roles"`
	jwt.RegisteredClaims
}

// VoterClaims 投票者 JWT 令牌
type VoterClaims struct {
	ID      uint64    `json:"id"`
	VoteID  uuid.UUID `json:"voteId"`
	IsVoted bool      `json:"isVoted"`
	jwt.RegisteredClaims
}

// GenUserToken 生成用戶 JWT 令牌
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

// GenVoterToken 生成投票者 JWT 令牌
func GenVoterToken(Id uint64, voteId uuid.UUID, isVoted bool) (string, string, error) {
	accessClaims := VoterClaims{
		ID:      Id,
		VoteID:  voteId,
		IsVoted: isVoted,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}

	return GenToken(accessClaims)
}

// GenToken 生成 JWT 令牌
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

// ParseUserToken 解析使用者 JWT 令牌
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

// ParseVoterToken 解析投票者 JWT 令牌
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

// parseToken 解析 JWT 令牌的通用函數
func parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
}

// ParseRefreshToken Parse and validate refresh token
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

// JWTAuthMiddleware 可選的 JWT middleware (適用於 GraphQL)
// 不會阻止請求,但會嘗試解析並設置所有可用的 token claims
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 嘗試從 Header 取得 token
		var headerToken string
		if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
			if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
				headerToken = parts[1]
			}
		}

		// 嘗試解析 user-token (從 Header 或 Cookie)
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

		// 嘗試解析 voter-token (從 Cookie,因為 Header 只能有一個)
		if voterToken, err := c.Cookie("voter-token"); err == nil && voterToken != "" {
			if voterClaims, err := ParseVoterToken(voterToken); err == nil {
				c.Set("voterId", voterClaims.ID)
				c.Set("voterVoteId", voterClaims.VoteID)
				c.Set("voterIsVoted", voterClaims.IsVoted)
				c.Set("hasVoterToken", true)
			}
		}

		c.Next()
	}
}

// RequireUserToken 在 GraphQL resolver 中驗證 user token
// 返回 UserClaims 或錯誤
func RequireUserToken(c *gin.Context) (*UserClaims, error) {
	// 先檢查是否已經由 middleware 設置
	if hasToken, exists := c.Get("hasUserToken"); exists && hasToken.(bool) {
		return &UserClaims{
			ID:      c.GetUint64("userId"),
			Account: c.GetString("userAccount"),
			Roles:   c.GetStringSlice("userRoles"),
		}, nil
	}

	// 否則嘗試手動解析
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

// RequireVoterToken 在 GraphQL resolver 中驗證 voter token
// 返回 VoterClaims 或錯誤
func RequireVoterToken(c *gin.Context) (*VoterClaims, error) {
	// 先檢查是否已經由 middleware 設置
	if hasToken, exists := c.Get("hasVoterToken"); exists && hasToken.(bool) {
		voteId, _ := c.Get("voterVoteId")
		return &VoterClaims{
			ID:      c.GetUint64("voterId"),
			VoteID:  voteId.(uuid.UUID),
			IsVoted: c.GetBool("voterIsVoted"),
		}, nil
	}

	// 否則嘗試手動解析
	var tokenString string
	if voterToken, err := c.Cookie("voter-token"); err == nil && voterToken != "" {
		tokenString = voterToken
	}

	if tokenString == "" {
		return nil, errors.New("voter token not found")
	}

	return ParseVoterToken(tokenString)
}

// GetOptionalUserToken 在 GraphQL resolver 中嘗試取得 user token (不強制)
// 返回 UserClaims 或 nil
func GetOptionalUserToken(c *gin.Context) *UserClaims {
	claims, _ := RequireUserToken(c)
	return claims
}

// GetOptionalVoterToken 在 GraphQL resolver 中嘗試取得 voter token (不強制)
// 返回 VoterClaims 或 nil
func GetOptionalVoterToken(c *gin.Context) *VoterClaims {
	claims, _ := RequireVoterToken(c)
	return claims
}
