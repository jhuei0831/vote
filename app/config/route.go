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
	// RBAC
	r.GET("/rbac/init", 
		middleware.JWTAuthMiddleware(),
		controller.NewRbacController().Initial,
	)
	// Anon
	r.POST("/v1/anon/login", controller.NewAnonController().AnonLogin)
	// User
	posts := r.Group("/v1/user")
	{
		posts.POST("/login", controller.NewUserController().Login)
		posts.POST("/check-auth",
			middleware.JWTAuthMiddleware(),
			controller.NewUserController().CheckAuth,
		)
		posts.POST("/refresh-token",
			middleware.JWTAuthMiddleware(),
			controller.NewUserController().RefreshToken,
		)
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
		votes.DELETE("/",
			middleware.RoleMiddleware("vote", "delete"),
			controller.NewVoteController().DeleteVote,
		)
	}

	// Question
	questions := r.Group("/v1/question", middleware.JWTAuthMiddleware())
	{
		questions.POST("/create",
			middleware.RoleMiddleware("question", "create"),
			controller.NewQuestionController().CreateQuestion,
		)
		questions.GET("/:id",
			middleware.RoleMiddleware("question", "read"),
			controller.NewQuestionController().SelectOneQuestion,
		)
		questions.GET("/list/:vote_id",
			middleware.RoleMiddleware("question", "read"),
			controller.NewQuestionController().SelectAllQuestions,
		)
		// questions.PUT("/:id",
		// 	middleware.RoleMiddleware("question", "update"),
		// 	controller.NewQuestionController().UpdateQuestion,
		// )
		// questions.DELETE("/",
		// 	middleware.RoleMiddleware("question", "delete"),
		// 	controller.NewQuestionController().DeleteQuestion,
		// )
	}

	// Candidate
	candidates := r.Group("/v1/candidate", middleware.JWTAuthMiddleware())
	{
		candidates.POST("/create",
			middleware.RoleMiddleware("candidate", "create"),
			controller.NewCandidateController().CreateCandidate,
		)
		candidates.GET("/:id",
			middleware.RoleMiddleware("candidate", "read"),
			controller.NewCandidateController().SelectOneCandidate,
		)
		candidates.GET("/list/:vote_id",
			middleware.RoleMiddleware("candidate", "read"),
			controller.NewCandidateController().SelectAllCandidates,
		)
		// candidates.PUT("/:id",
		// 	middleware.RoleMiddleware("candidate", "update"),
		// 	controller.NewCandidateController().UpdateCandidate,
		// )
		// candidates.DELETE("/",
		// 	middleware.RoleMiddleware("candidate", "delete"),
		// 	controller.NewCandidateController().DeleteCandidate,
		// )
	}

	// Password
	passwords := r.Group("/v1/password", middleware.JWTAuthMiddleware())
	{
		passwords.POST("/create",
			middleware.RoleMiddleware("password", "create"),
			controller.NewPasswordController().CreatePassword,
		)
		passwords.POST("/decrypt",
			middleware.RoleMiddleware("password", "read"),
			controller.NewPasswordController().DecryptPassword,
		)
		passwords.GET("/list/:vote_id",
			middleware.RoleMiddleware("password", "read"),
			controller.NewPasswordController().SelectAllPasswords,
		)
	}
}
