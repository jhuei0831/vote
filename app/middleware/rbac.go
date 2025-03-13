package middleware

// import "github.com/gin-gonic/gin"

// func rbac() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Check if the user has the right to access the resource
// 		// If the user has the right to access the resource, call ctx.Next()
// 		// If the user does not have the right to access the resource, return an error message
// 		tokenString := c.GetHeader("Authorization")
//         token, _ := jwt.Parse(tokenString)
//         claims, _ := token.Claims.(jwt.MapClaims)
//         role := claims["role"].(string)

//         if role != RoleAdmin {
//             c.AbortWithStatus(http.StatusForbidden)
//             return
//         }
//         c.Next()
// 	}
// }
