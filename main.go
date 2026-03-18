package main

import (
    "os"

    "vote/app/config"
    "vote/app/database"
    "vote/app/middleware"
    "vote/internal/grpcserver"

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
    if err := godotenv.Load(); err != nil {
        panic(err)
    }

    // Initialize database
    dbConfig := database.DbConfig()
    _, err := database.Initialize(dbConfig)
    if err != nil {
        panic(err)
    }
    
    // Initialize RBAC
    _, _, err = database.Rbac()
    if err != nil {
        panic(err)
    }

    // Start gRPC server in background
    go func() {
        grpcserver.StartGRPCServer("50051")
    }()

    // Start HTTP server (blocking)
    port := os.Getenv("PORT")
    server := SetRouter()
    if err := server.Run(":" + port); err != nil {
        panic(err)
    }
}

func SetRouter() *gin.Engine {
    server := gin.Default()
    server.Use(middleware.GinContextToContextMiddleware())
    server.Use(middleware.CORSMiddleware())
    server.Use(middleware.LoggerToFile())
    config.Routes(server, config.RedisStore())
    
    return server
}