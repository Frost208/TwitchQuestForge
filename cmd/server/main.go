package main

import (
        "log"
        "os"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/handlers"

        "github.com/gin-gonic/gin"
        "github.com/joho/godotenv"
)

func main() {
        log.Println("=== Twitch RPG Server Starting ===")
        
        // Load environment variables
        if err := godotenv.Load(); err != nil {
                log.Println("No .env file found, using environment variables")
        }

        log.Println("Attempting database connection...")
        // Connect to database
        if err := database.Connect(); err != nil {
                log.Printf("Warning: Failed to connect to database: %v", err)
                log.Println("Server will start without database connection for testing")
        } else {
                defer database.Close()
                log.Println("Database connected successfully")
        }

        log.Println("Setting up HTTP server...")
        // Set Gin mode based on environment
        if os.Getenv("GIN_MODE") == "release" {
                gin.SetMode(gin.ReleaseMode)
        }

        // Create Gin router
        router := gin.Default()

        // Add CORS middleware
        router.Use(func(c *gin.Context) {
                c.Header("Access-Control-Allow-Origin", "*")
                c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

                if c.Request.Method == "OPTIONS" {
                        c.AbortWithStatus(204)
                        return
                }

                c.Next()
        })

        log.Println("Registering API routes...")
        // Register API routes
        handlers.RegisterRoutes(router)

        // Start server
        port := os.Getenv("SERVER_PORT")
        if port == "" {
                port = "8080"
        }

        log.Printf("Starting Twitch RPG server on port %s", port)
        log.Println("Server is ready to accept connections!")
        if err := router.Run(":" + port); err != nil {
                log.Fatal("Failed to start server:", err)
        }
}