package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dirgocs/supabase-self-hosted-mcp/config"
	"github.com/dirgocs/supabase-self-hosted-mcp/controllers"
	"github.com/dirgocs/supabase-self-hosted-mcp/supabase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// Initialize extended Supabase client with Functions support
	supabaseClient := supabase.CreateClientExtended(cfg.Supabase.URL, cfg.Supabase.Key)
	
	// The headers are already set up in the CreateClientExtended function
	// No need to manually set them here

	// Set up Gin router
	router := gin.Default()

	// Configure trusted proxies
	// For Cloudflare, we need to trust their IP ranges
	// See: https://www.cloudflare.com/ips/
	// For simplicity in a Docker/Coolify environment, we can trust the local subnet
	// and Cloudflare's proxy
	trustedProxies := []string{
		"10.0.0.0/8",     // Private network range often used in Docker
		"172.16.0.0/12",  // Private network range often used in Docker
		"192.168.0.0/16", // Private network range
		// Cloudflare IPv4 ranges - you may want to keep these updated
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
	}
	router.SetTrustedProxies(trustedProxies)

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
	edgeFunctionsController := controllers.NewEdgeFunctionsController(supabaseClient)
	router.POST("/v1/get_edge_functions", edgeFunctionsController.GetEdgeFunctions)
	router.POST("/v1/create_edge_function", edgeFunctionsController.CreateEdgeFunction)
	router.POST("/v1/update_edge_function", edgeFunctionsController.UpdateEdgeFunction)
	router.POST("/v1/delete_edge_function", edgeFunctionsController.DeleteEdgeFunction)
	router.POST("/v1/deploy_edge_function", edgeFunctionsController.DeployEdgeFunction)

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
