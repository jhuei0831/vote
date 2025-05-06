package middleware

import (
	"errors"
	"net/http"
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
	ID			uint64 			`json:"id"`
	VoteID  uuid.UUID 	`json:"voteId"`
	IsVoted   bool			`json:"isVoted"`
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

// JWTAuthMiddleware Middleware of JWT
func JWTAuthMiddleware(isUser bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenType := "voter"
		if isUser {
			tokenType = "user"
		}

		// 從 Header 或 Cookie 取得 token
		tokenString := extractToken(c, tokenType)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Authorization token not found in Header or Cookie",
			})
			c.Abort()
			return
		}

		// 驗證並解析 token
		if err := validateAndSetClaims(c, tokenString, isUser); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken 從 Header 或 Cookie 取得 token
func extractToken(c *gin.Context, tokenType string) string {
	// 先從 Header 取得
	if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 若 Header 沒有則從 Cookie 取得
	if tokenCookie, err := c.Cookie(tokenType + "-token"); err == nil && tokenCookie != "" {
		return tokenCookie
	}

	return ""
}

// validateAndSetClaims 驗證 token 並設置 claims
func validateAndSetClaims(c *gin.Context, tokenString string, isUser bool) error {
	if isUser {
		mc, err := ParseUserToken(tokenString)
		if err != nil {
			return err
		}
		c.Set("id", mc.ID)
		c.Set("account", mc.Account)
		c.Set("roles", mc.Roles)
	} else {
		mc, err := ParseVoterToken(tokenString)
		if err != nil {
			return err
		}
		c.Set("id", mc.ID)
		c.Set("voteId", mc.VoteID)
	}
	return nil
}
