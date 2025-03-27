package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var RefreshSecretKey = []byte(os.Getenv("JWT_REFRESH_SECRET_KEY"))

const TokenExpireDuration = time.Hour * 2
const RefreshTokenExpireDuration = time.Hour * 24 * 7 // Refresh Token 有效期 7 天

type MyClaims struct {
	ID      uint64   `json:"id"`
	Account string   `json:"account"`
	Roles   []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenToken Create a new access and refresh token
func GenToken(Id uint64, account string, roles []string) (string, string, error) {
	// Access Token
	accessClaims := MyClaims{
		ID:      Id,
		Account: account,
		Roles:   roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    os.Getenv("APP_NAME"),
		},
	}
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

// ParseToken Parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// Valid token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		var tokenString string

		// Check token in Header.Authorization field
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// If no valid token in header, check in cookies
		if tokenString == "" {
			tokenCookie, err := c.Cookie("token")
			if err != nil || tokenCookie == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": -1,
					"msg":  "Authorization token not found in Header or Cookie",
				})
				c.Abort()
				return
			}
			tokenString = tokenCookie
		}

		// Parse and validate the token
		mc, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Invalid Token.",
			})
			c.Abort()
			return
		}

		// Store Account info into Context
		c.Set("id", mc.ID)
		c.Set("account", mc.Account)
		c.Set("roles", mc.Roles)

		// Proceed to the next middleware or handler
		c.Next()
	}
}
