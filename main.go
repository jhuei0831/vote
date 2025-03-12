package main

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"vote/app/config"
	"vote/app/database"
	"vote/app/model"
	"vote/app/middleware"
	"os"
)

// @title Gin swagger
// @version 1.0
// @description Gin swagger

// @contact.name Flynn Sun

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
	db, err := database.Initialize(dbConfig)
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}

	server := gin.Default()
	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB: 0,
	}))
	config.Routes(server, redisStore)
	config.Swagger()
	server.Use(middleware.LoggerToFile())
	server.Use(middleware.CORSMiddleware())

	return server
}
