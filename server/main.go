package main

import (
	"log"
	"path/filepath"
	"veritas-server/api"
	"veritas-server/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from project root (parent directory)
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		// Try current directory as fallback
		if err2 := godotenv.Load(); err2 != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Initialize Database
	db.Init()

	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	// API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/models", api.GetModels)
		apiGroup.POST("/chat", api.Chat)
		apiGroup.GET("/conversations", api.GetConversations)
		apiGroup.GET("/conversations/:id", api.GetConversation)
		apiGroup.POST("/conversations", api.CreateConversation)

		// Model configuration endpoints
		apiGroup.POST("/model-configs", api.CreateModelConfig)
		apiGroup.GET("/model-configs", api.GetModelConfigs)
		apiGroup.GET("/model-configs/:id", api.GetModelConfig)
		apiGroup.PUT("/model-configs/:id", api.UpdateModelConfig)
		apiGroup.DELETE("/model-configs/:id", api.DeleteModelConfig)
		apiGroup.POST("/model-configs/test", api.TestModelConfig)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
