package middleware

import (
	"net/http"
	"vote/app/database"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the right to access the resource
		// If the user has the right to access the resource, call ctx.Next()
		// If the user does not have the right to access the resource, return an error message
		account, exists := c.Get("account")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Account not found"})
			return
		}

		ok, err := database.Enforcer.Enforce(account.(string), obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

        c.Next()
	}
}
