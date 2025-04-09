package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dirgocs/supabase-self-hosted-mcp/config"
	"github.com/dirgocs/supabase-self-hosted-mcp/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file, using default values")
	}

	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize Supabase client
	supabaseClient := supabase.CreateClient(cfg.Supabase.URL, cfg.Supabase.Key)
	
	// Set up the client with additional options
	// Add the service role key to the Auth header for admin operations
	supabaseClient.Auth.Header.Add("apikey", cfg.Supabase.Key)
	// Also add it as a Bearer token for services that require it
	supabaseClient.Auth.Header.Add("Authorization", "Bearer "+cfg.Supabase.Key)

	// Set up Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(cors.Default())

	// MCP Server info endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Supabase Self-Hosted MCP Server",
			"version":     "1.0.0",
			"description": "MCP Server para comunicação com Supabase Self-Hosted",
			"status":      "running",
		})
	})
	
	// Health check endpoint for Docker
	router.GET("/health", func(c *gin.Context) {
		// Check if Supabase is accessible
		supabaseStatus := "unknown"
		
		// Try to make a simple request to Supabase
		resp, err := http.Get(cfg.Supabase.URL)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			supabaseStatus = "connected"
		} else {
			supabaseStatus = "disconnected"
		}
		
		c.JSON(200, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"supabase": supabaseStatus,
		})
	})

	// MCP specification endpoint
	router.GET("/v1/specification", controllers.GetMCPSpecification)

	// Register database endpoints
	dbController := controllers.NewDatabaseController(supabaseClient)
	router.POST("/v1/execute_query", dbController.ExecuteQuery)
	router.POST("/v1/get_database_schema", dbController.GetDatabaseSchema)
	router.POST("/v1/create_schema", dbController.CreateSchema)
	router.POST("/v1/delete_schema", dbController.DeleteSchema)
	router.POST("/v1/get_rls_policies", dbController.GetRLSPolicies)
	router.POST("/v1/create_rls_policy", dbController.CreateRLSPolicy)
	router.POST("/v1/update_rls_policy", dbController.UpdateRLSPolicy)
	router.POST("/v1/delete_rls_policy", dbController.DeleteRLSPolicy)

	// Register table endpoints
	tableController := controllers.NewTableController(supabaseClient)
	router.POST("/v1/query_table", tableController.QueryTable)
	router.POST("/v1/generate_types", tableController.GenerateTypes)
	router.POST("/v1/list_tables", tableController.ListTables)
	router.POST("/v1/create_table", tableController.CreateTable)
	router.POST("/v1/alter_table", tableController.AlterTable)
	router.POST("/v1/drop_table", tableController.DropTable)

	// Register storage endpoints
	storageController := controllers.NewStorageController(supabaseClient)
	router.POST("/v1/get_buckets", storageController.GetBuckets)
	router.POST("/v1/create_bucket", storageController.CreateBucket)
	router.POST("/v1/update_bucket", storageController.UpdateBucket)
	router.POST("/v1/delete_bucket", storageController.DeleteBucket)
	router.POST("/v1/get_bucket_policies", storageController.GetBucketPolicies)
	router.POST("/v1/create_bucket_policy", storageController.CreateBucketPolicy)
	router.POST("/v1/update_bucket_policy", storageController.UpdateBucketPolicy)
	router.POST("/v1/delete_bucket_policy", storageController.DeleteBucketPolicy)

	// Register edge function endpoints
	edgeFunctionController := controllers.NewEdgeFunctionController(supabaseClient)
	router.POST("/v1/get_edge_functions", edgeFunctionController.GetEdgeFunctions)
	router.POST("/v1/create_edge_function", edgeFunctionController.CreateEdgeFunction)
	router.POST("/v1/update_edge_function", edgeFunctionController.UpdateEdgeFunction)
	router.POST("/v1/delete_edge_function", edgeFunctionController.DeleteEdgeFunction)
	router.POST("/v1/deploy_edge_function", edgeFunctionController.DeployEdgeFunction)

	// Start the server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Supabase Self-Hosted MCP Server running on port %d", cfg.Server.Port)
	log.Printf("Server URL: http://localhost%s", port)
	log.Printf("MCP Specification URL: http://localhost%s/v1/specification", port)
	log.Printf("Supabase URL: %s", cfg.Supabase.URL)

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
