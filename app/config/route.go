package config

import (
	"os"
	// "time"

	"vote/app/controller"
	"vote/app/middleware"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	// cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

func Routes(r *gin.Engine, m *persist.RedisStore) {
	// Graphql
	r.POST("/query", middleware.JWTAuthMiddleware(), graphqlHandler())
	r.GET("/", playgroundHandler())

	// Restful API
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health check: PORT " + os.Getenv("APP_PORT"),
		})
		utils.Logger().WithFields(logrus.Fields{
			"name": os.Getenv("APP_NAME"),
		}).Info("Health Check", "Info")
	})

	// RBAC
	r.GET("/rbac/init",
		middleware.JWTAuthMiddleware(),
		controller.NewRbacController().Initial,
	)
	
	// Voter
	voter := r.Group("/v1/voter")
	{
		// voter.POST("/login", controller.NewVoterController().VoterLogin)
		voter.POST("/logout",
			middleware.JWTAuthMiddleware(),
			controller.NewVoterController().Logout,
		)
		// voter.POST("/check-auth",
		// 	middleware.JWTAuthMiddleware(),
		// 	controller.NewVoterController().CheckAuth,
		// )
	}

	// User
	users := r.Group("/v1/user")
	{
		users.POST("/login", controller.NewUserController().Login)
		users.POST("/logout",
			middleware.JWTAuthMiddleware(),
			controller.NewUserController().Logout,
		)
		users.POST("/check-auth",
			middleware.JWTAuthMiddleware(),
			controller.NewUserController().CheckAuth,
		)
		users.POST("/refresh-token",
			middleware.JWTAuthMiddleware(),
			controller.NewUserController().RefreshToken,
		)
		users.POST("/create",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware("user", "create"),
			controller.NewUserController().CreateUser,
		)
		users.GET("/me",
			middleware.JWTAuthMiddleware(),
			// middleware.RoleMiddleware("user", "read"),
			// cache.CacheByRequestURI(m, 2*time.Hour),
			controller.NewUserController().GetUser,
		)
	}
}
