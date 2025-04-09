package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

// EdgeFunctionController handles edge function-related operations
type EdgeFunctionController struct {
	supabase *supabase.Client
}

// NewEdgeFunctionController creates a new edge function controller
func NewEdgeFunctionController(supabase *supabase.Client) *EdgeFunctionController {
	return &EdgeFunctionController{
		supabase: supabase,
	}
}

// GetEdgeFunctionsRequest represents the request body for getting edge functions
type GetEdgeFunctionsRequest struct {
	Name string `json:"name"`
}

// GetEdgeFunctions gets edge functions
func (efc *EdgeFunctionController) GetEdgeFunctions(c *gin.Context) {
	var req GetEdgeFunctionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For self-hosted Supabase, we need to check where edge functions are stored
	// This may vary depending on the setup

	// Option 1: If edge functions are stored in a specific table
	query := "SELECT * FROM edge_functions"

	if req.Name != "" {
		query += fmt.Sprintf(" WHERE name = '%s'", req.Name)
	}

	var result []map[string]interface{}
	err := efc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": query,
	}, &result)

	if err == nil && result != nil {
		c.JSON(http.StatusOK, result)
		return
	}

	// Option 2: For self-hosted implementations, an alternative method may be needed
	// As a fallback, return information that a specific implementation is needed
	c.JSON(http.StatusOK, gin.H{
		"message":              "Edge functions management requires a specific implementation for your self-hosted Supabase setup",
		"implementation_needed": true,
		"functions":            []interface{}{},
	})
}

// ImportMap represents an import map for edge functions
type ImportMap map[string]string

// CreateEdgeFunctionRequest represents the request body for creating an edge function
type CreateEdgeFunctionRequest struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	VerifyJWT bool      `json:"verify_jwt"`
	ImportMap ImportMap `json:"import_map"`
}

// CreateEdgeFunction creates a new edge function
func (efc *EdgeFunctionController) CreateEdgeFunction(c *gin.Context) {
	var req CreateEdgeFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// For self-hosted Supabase, edge function creation implementation

	// Option 1: If there's a table for storing edge functions
	importMapJSON, err := json.Marshal(req.ImportMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import map format"})
		return
	}

	// Escape single quotes in code and import map
	code := strings.ReplaceAll(req.Code, "'", "''")
	importMapStr := strings.ReplaceAll(string(importMapJSON), "'", "''")

	insertQuery := fmt.Sprintf(`
		INSERT INTO edge_functions (name, code, verify_jwt, import_map, created_at, updated_at)
		VALUES (
			'%s',
			'%s',
			%t,
			'%s',
			NOW(),
			NOW()
		)
	`, req.Name, code, req.VerifyJWT, importMapStr)

	var insertResult interface{}
	err = efc.supabase.DB.RPC("execute_sql", map[string]interface{}{
		"query": insertQuery,
	}, &insertResult)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Edge function '%s' created successfully", req.Name),
		})
		return
	}

	// Option 2: Specific implementation for self-hosted
	c.JSON(http.StatusOK, gin.H{
		"message":              "Edge function creation requires a specific implementation for your self-hosted Supabase setup",
		"implementation_needed": true,
		"function_name":        req.Name,
	})
}

// UpdateEdgeFunctionRequest represents the request body for updating an edge function
type UpdateEdgeFunctionRequest struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	VerifyJWT *bool     `json:"verify_jwt"`
	ImportMap ImportMap `json:"import_map"`
}

// UpdateEdgeFunction updates an existing edge function
func (efc *EdgeFunctionController) UpdateEdgeFunction(c *gin.Context) {
	var req UpdateEdgeFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// For self-hosted Supabase, specific implementation for update

	// Option 1: If there's a table for storing edge functions
	// Escape single quotes in code
	code := strings.ReplaceAll(req.Code, "'", "''")

	updateQuery := fmt.Sprintf(`
		UPDATE edge_functions 
		SET code = '%s'
	`, code)

	if req.VerifyJWT != nil {
		updateQuery += fmt.Sprintf(`, verify_jwt = %t`, *req.VerifyJWT)
	}

	if req.ImportMap != nil {
		importMapJSON, err := json.Marshal(req.ImportMap)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import map format"})
			return
		}
		// Escape single quotes in import map
		importMapStr := strings.ReplaceAll(string(importMapJSON), "'", "''")
		updateQuery += fmt.Sprintf(`, import_map = '%s'`, importMapStr)
	}

	updateQuery += fmt.Sprintf(`, updated_at = NOW() WHERE name = '%s'`, req.Name)

	var updateResult interface{}
	err := efc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": updateQuery,
	}, &updateResult)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Edge function '%s' updated successfully", req.Name),
		})
		return
	}

	// Option 2: Specific implementation for self-hosted
	c.JSON(http.StatusOK, gin.H{
		"message":              "Edge function update requires a specific implementation for your self-hosted Supabase setup",
		"implementation_needed": true,
		"function_name":        req.Name,
	})
}

// DeleteEdgeFunctionRequest represents the request body for deleting an edge function
type DeleteEdgeFunctionRequest struct {
	Name string `json:"name"`
}

// DeleteEdgeFunction deletes an edge function
func (efc *EdgeFunctionController) DeleteEdgeFunction(c *gin.Context) {
	var req DeleteEdgeFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function name is required"})
		return
	}

	// Implementation for self-hosted Supabase

	// Option 1: If there's a table for storing edge functions
	deleteQuery := fmt.Sprintf(`
		DELETE FROM edge_functions 
		WHERE name = '%s'
	`, req.Name)

	var deleteResult interface{}
	err := efc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": deleteQuery,
	}, &deleteResult)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Edge function '%s' deleted successfully", req.Name),
		})
		return
	}

	// Option 2: Specific implementation for self-hosted
	c.JSON(http.StatusOK, gin.H{
		"message":              "Edge function deletion requires a specific implementation for your self-hosted Supabase setup",
		"implementation_needed": true,
		"function_name":        req.Name,
	})
}

// DeployEdgeFunctionRequest represents the request body for deploying an edge function
type DeployEdgeFunctionRequest struct {
	Name string `json:"name"`
}

// DeployEdgeFunction deploys an edge function
func (efc *EdgeFunctionController) DeployEdgeFunction(c *gin.Context) {
	var req DeployEdgeFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function name is required"})
		return
	}

	// For self-hosted Supabase, the deployment process may be specific
	// Basic implementation that will need to be adapted

	c.JSON(http.StatusOK, gin.H{
		"message":              "Edge function deployment requires a specific implementation for your self-hosted Supabase setup",
		"implementation_needed": true,
		"function_name":        req.Name,
	})
}
