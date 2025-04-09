package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default values")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Set up Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.Default())

	// MCP Server info endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Supabase Self-Hosted MCP Server",
			"version":     "1.0.0",
			"description": "MCP Server for communication with Supabase Self-Hosted",
			"status":      "running",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Start the server
	serverPort := fmt.Sprintf(":%s", port)
	log.Printf("Supabase Self-Hosted MCP Server running on port %s", port)
	log.Printf("Server URL: http://localhost%s", serverPort)

	if err := router.Run(serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
