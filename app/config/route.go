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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(r *gin.Engine, m *persist.RedisStore) {
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health check: PORT " + os.Getenv("PORT"),
		})
		utils.Logger().WithFields(logrus.Fields{
			"name": os.Getenv("APP_NAME"),
		}).Info("Health Check", "Info")
	})
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// User
	posts := r.Group("/v1/user")
	{
		posts.POST("/login", controller.NewUserController().AuthHandler)
		posts.POST("/create", 
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware("user", "create"),
			controller.NewUserController().CreateUser,
		)
		posts.GET("/:id",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware("user", "read"),
			// cache.CacheByRequestURI(m, 2*time.Hour),
			controller.NewUserController().GetUser,
		)
	}

	// Vote
	votes := r.Group("/v1/vote", middleware.JWTAuthMiddleware())
	{
		votes.POST("/create",
			middleware.RoleMiddleware("vote", "create"),
			controller.NewVoteController().CreateVote,
		)
		votes.GET("/:id",
			middleware.RoleMiddleware("vote", "read"),
			controller.NewVoteController().SelectOneVote,
		)
		votes.GET("/list",
			middleware.RoleMiddleware("vote", "read"),
			controller.NewVoteController().SelectAllVotes,
		)
		votes.PUT("/:id",
			middleware.RoleMiddleware("vote", "update"),
			controller.NewVoteController().UpdateVote,
		)
	}
}
