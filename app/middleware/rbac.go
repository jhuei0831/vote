package middleware

import (
	"net/http"
	"strconv"
	"vote/app/database"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the right to access the resource
		// If the user has the right to access the resource, call ctx.Next()
		// If the user does not have the right to access the resource, return an error message
		id, exists := c.Get("id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Account not found"})
			return
		}
		userID := id.(uint64)
		userId := strconv.FormatUint(userID, 10)
		ok, err := database.Enforcer.Enforce(userId, obj, act)
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
