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
	// Graphql
	r.POST("/query", middleware.JWTAuthMiddleware(true), graphqlHandler())
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
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// RBAC
	r.GET("/rbac/init",
		middleware.JWTAuthMiddleware(true),
		controller.NewRbacController().Initial,
	)
	// Voter
	r.POST("/v1/voter/login", controller.NewVoterController().VoterLogin)
	r.POST("/v1/voter/logout",
		middleware.JWTAuthMiddleware(false),
		controller.NewVoterController().Logout,
	)
	r.POST("/v1/voter/check-auth",
		middleware.JWTAuthMiddleware(false),
		controller.NewVoterController().CheckAuth,
	)
	// r.GET("/v1/voter/questions",
	// 	middleware.JWTAuthMiddleware(false),
	// 	controller.NewQuestionController().SelectVoterQuestions,
	// )
	r.POST("/v1/voter/ballot/create",
		middleware.JWTAuthMiddleware(false),
		controller.NewBallotController().CreateBallots,
	)
	// User
	posts := r.Group("/v1/user")
	{
		posts.POST("/login", controller.NewUserController().Login)
		posts.POST("/logout",
			middleware.JWTAuthMiddleware(true),
			controller.NewUserController().Logout,
		)
		posts.POST("/check-auth",
			middleware.JWTAuthMiddleware(true),
			controller.NewUserController().CheckAuth,
		)
		posts.POST("/refresh-token",
			middleware.JWTAuthMiddleware(true),
			controller.NewUserController().RefreshToken,
		)
		posts.POST("/create",
			middleware.JWTAuthMiddleware(true),
			middleware.RoleMiddleware("user", "create"),
			controller.NewUserController().CreateUser,
		)
		posts.GET("/:id",
			middleware.JWTAuthMiddleware(true),
			middleware.RoleMiddleware("user", "read"),
			// cache.CacheByRequestURI(m, 2*time.Hour),
			controller.NewUserController().GetUser,
		)
	}

	// Vote
	r.GET("/v1/vote/:id", controller.NewVoteController().GetVote)
	votes := r.Group("/v1/vote", middleware.JWTAuthMiddleware(true))
	{
		votes.POST("/create",
			middleware.RoleMiddleware("vote", "create"),
			controller.NewVoteController().CreateVote,
		)
		// votes.GET("/:id",
		// 	middleware.RoleMiddleware("vote", "read"),
		// 	controller.NewVoteController().GetVote,
		// )
		// votes.GET("/list",
		// 	middleware.RoleMiddleware("vote", "read"),
		// 	controller.NewVoteController().GetVotes,
		// )
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
	questions := r.Group("/v1/question", middleware.JWTAuthMiddleware(true))
	{
		questions.POST("/create",
			middleware.RoleMiddleware("question", "create"),
			controller.NewQuestionController().CreateQuestion,
		)
		questions.GET("/:id",
			middleware.RoleMiddleware("question", "read"),
			controller.NewQuestionController().GetQuestion,
		)
		// questions.GET("/list/:vote_id",
		// 	middleware.RoleMiddleware("question", "read"),
		// 	controller.NewQuestionController().GetQuestions,
		// )
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
	candidates := r.Group("/v1/candidate", middleware.JWTAuthMiddleware(true))
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
	passwords := r.Group("/v1/password", middleware.JWTAuthMiddleware(true))
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
		passwords.PUT("/update-status",
			middleware.RoleMiddleware("password", "update"),
			controller.NewPasswordController().UpdatePasswordStatus,
		)
	}
}
