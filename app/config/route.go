package config

import (
	"os"
	"vote/app/controller"
	"vote/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"time"
)

func Routes(r *gin.Engine, m *persist.RedisStore) {
	r.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health check: PORT 2 " + os.Getenv("PORT"),
		})
		middleware.Logger().WithFields(logrus.Fields{
			"name": "Flynn Sun",
		}).Info("Health Check", "Info")
	})
	// Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// User
    posts := r.Group("/v1/users")
    {
        posts.POST("/", controller.NewUsersController().CreateUser)
        posts.POST("/login", controller.QueryUsersController().AuthHandler)
        posts.GET("/:id", middleware.JWTAuthMiddleware(), cache.CacheByRequestURI(m, 2*time.Hour), controller.QueryUsersController().GetUser)
    }
}