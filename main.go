package main

import (
	"os"

	"vote/app/config"
	"vote/app/database"
	"vote/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Gin swagger
// @version 1.0
// @description Gin swagger

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}

	port := os.Getenv("PORT")
	server := SetRouter()
	err := server.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func SetRouter() *gin.Engine {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Initialize database
	dbConfig := os.Getenv("DB_CONFIG")
	_, err := database.Initialize(dbConfig)
	if err != nil {
		panic(err)
	}
	
	// Initialize RBAC
	_, _, err = database.Rbac()
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.Use(middleware.LoggerToFile())
	config.Routes(server, config.RedisStore())
	config.Swagger()
	server.Use(middleware.CORSMiddleware())

	return server
}
