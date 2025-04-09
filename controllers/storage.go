package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

// StorageController handles storage-related operations
type StorageController struct {
	supabase *supabase.Client
}

// NewStorageController creates a new storage controller
func NewStorageController(supabase *supabase.Client) *StorageController {
	return &StorageController{
		supabase: supabase,
	}
}

// GetBucketsRequest represents the request body for getting storage buckets
type GetBucketsRequest struct {
	ID string `json:"id"`
}

// GetBuckets gets storage buckets
func (sc *StorageController) GetBuckets(c *gin.Context) {
	var req GetBucketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT * FROM storage.buckets`

	if req.ID != "" {
		query += fmt.Sprintf(` WHERE id = '%s'`, req.ID)
	}

	var result []map[string]interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": query,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateBucketRequest represents the request body for creating a storage bucket
type CreateBucketRequest struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Public           bool     `json:"public"`
	FileSizeLimit    *int64   `json:"file_size_limit"`
	AllowedMimeTypes []string `json:"allowed_mime_types"`
}

// CreateBucket creates a new storage bucket
func (sc *StorageController) CreateBucket(c *gin.Context) {
	var req CreateBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket ID is required"})
		return
	}

	// Try using Supabase API first
	// Note: The current supabase-go client doesn't support storage operations directly
	// We'll implement the SQL fallback for now

	// Set default name if not provided
	if req.Name == "" {
		req.Name = req.ID
	}

	// Convert allowed mime types to JSON if provided
	var allowedMimeTypesJSON string
	if len(req.AllowedMimeTypes) > 0 {
		mimeTypesBytes, err := json.Marshal(req.AllowedMimeTypes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MIME types format"})
			return
		}
		allowedMimeTypesJSON = string(mimeTypesBytes)
	}

	// Build SQL query
	fileSizeLimitStr := "NULL"
	if req.FileSizeLimit != nil {
		fileSizeLimitStr = fmt.Sprintf("%d", *req.FileSizeLimit)
	}

	allowedMimeTypesStr := "NULL"
	if allowedMimeTypesJSON != "" {
		allowedMimeTypesStr = fmt.Sprintf("'%s'::jsonb", allowedMimeTypesJSON)
	}

	sql := fmt.Sprintf(`
		INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types, created_at, updated_at)
		VALUES (
			'%s',
			'%s',
			%t,
			%s,
			%s,
			NOW(),
			NOW()
		)
	`, req.ID, req.Name, req.Public, fileSizeLimitStr, allowedMimeTypesStr)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Bucket '%s' created successfully", req.ID),
		"method":  "sql",
	})
}

// UpdateBucketRequest represents the request body for updating a storage bucket
type UpdateBucketRequest struct {
	ID               string   `json:"id"`
	Public           *bool    `json:"public"`
	FileSizeLimit    *int64   `json:"file_size_limit"`
	AllowedMimeTypes []string `json:"allowed_mime_types"`
}

// UpdateBucket updates a storage bucket
func (sc *StorageController) UpdateBucket(c *gin.Context) {
	var req UpdateBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket ID is required"})
		return
	}

	// Build SQL query
	sql := `UPDATE storage.buckets SET updated_at = NOW()`

	if req.Public != nil {
		sql += fmt.Sprintf(`, public = %t`, *req.Public)
	}

	if req.FileSizeLimit != nil {
		sql += fmt.Sprintf(`, file_size_limit = %d`, *req.FileSizeLimit)
	} else {
		// If explicitly set to nil in the request
		if c.Request.ContentLength > 0 && strings.Contains(c.Request.Body.(*gin.Request).Body, "file_size_limit") {
			sql += `, file_size_limit = NULL`
		}
	}

	if req.AllowedMimeTypes != nil {
		if len(req.AllowedMimeTypes) > 0 {
			mimeTypesBytes, err := json.Marshal(req.AllowedMimeTypes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MIME types format"})
				return
			}
			sql += fmt.Sprintf(`, allowed_mime_types = '%s'::jsonb`, string(mimeTypesBytes))
		} else {
			sql += `, allowed_mime_types = NULL`
		}
	}

	sql += fmt.Sprintf(` WHERE id = '%s'`, req.ID)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Bucket '%s' updated successfully", req.ID),
		"method":  "sql",
	})
}

// DeleteBucketRequest represents the request body for deleting a storage bucket
type DeleteBucketRequest struct {
	ID string `json:"id"`
}

// DeleteBucket deletes a storage bucket
func (sc *StorageController) DeleteBucket(c *gin.Context) {
	var req DeleteBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket ID is required"})
		return
	}

	sql := fmt.Sprintf(`DELETE FROM storage.buckets WHERE id = '%s'`, req.ID)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Bucket '%s' deleted successfully", req.ID),
		"method":  "sql",
	})
}

// GetBucketPoliciesRequest represents the request body for getting bucket policies
type GetBucketPoliciesRequest struct {
	BucketID string `json:"bucket_id"`
}

// GetBucketPolicies gets policies for a storage bucket
func (sc *StorageController) GetBucketPolicies(c *gin.Context) {
	var req GetBucketPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BucketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket ID is required"})
		return
	}

	query := fmt.Sprintf(`
		SELECT 
			policyname AS name,
			CASE operation
				WHEN 10 THEN 'SELECT'
				WHEN 20 THEN 'INSERT'
				WHEN 40 THEN 'UPDATE'
				WHEN 80 THEN 'DELETE'
				ELSE 'UNKNOWN'
			END AS operation,
			definition,
			role
		FROM storage.policies
		WHERE name = '%s'
	`, req.BucketID)

	var result []map[string]interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": query,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateBucketPolicyRequest represents the request body for creating a bucket policy
type CreateBucketPolicyRequest struct {
	BucketID   string `json:"bucket_id"`
	Name       string `json:"name"`
	Operation  string `json:"operation"`
	Definition string `json:"definition"`
	Role       string `json:"role"`
}

// CreateBucketPolicy creates a policy for a storage bucket
func (sc *StorageController) CreateBucketPolicy(c *gin.Context) {
	var req CreateBucketPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BucketID == "" || req.Name == "" || req.Operation == "" || req.Definition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "authenticated"
	}

	// Map operation string to internal code
	opMap := map[string]int{
		"SELECT": 10,
		"INSERT": 20,
		"UPDATE": 40,
		"DELETE": 80,
	}

	opCode, ok := opMap[req.Operation]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation. Must be SELECT, INSERT, UPDATE, or DELETE"})
		return
	}

	// Escape single quotes in definition
	definition := strings.ReplaceAll(req.Definition, "'", "''")

	sql := fmt.Sprintf(`
		INSERT INTO storage.policies (name, bucket_id, operation, definition, role, created_at)
		VALUES (
			'%s',
			'%s',
			%d,
			'%s',
			'%s',
			NOW()
		)
	`, req.Name, req.BucketID, opCode, definition, req.Role)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Policy '%s' for bucket '%s' created successfully", req.Name, req.BucketID),
	})
}

// UpdateBucketPolicyRequest represents the request body for updating a bucket policy
type UpdateBucketPolicyRequest struct {
	BucketID   string `json:"bucket_id"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// UpdateBucketPolicy updates a policy for a storage bucket
func (sc *StorageController) UpdateBucketPolicy(c *gin.Context) {
	var req UpdateBucketPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BucketID == "" || req.Name == "" || req.Definition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Escape single quotes in definition
	definition := strings.ReplaceAll(req.Definition, "'", "''")

	sql := fmt.Sprintf(`
		UPDATE storage.policies
		SET 
			definition = '%s',
			updated_at = NOW()
		WHERE 
			bucket_id = '%s' AND 
			name = '%s'
	`, definition, req.BucketID, req.Name)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Policy '%s' for bucket '%s' updated successfully", req.Name, req.BucketID),
	})
}

// DeleteBucketPolicyRequest represents the request body for deleting a bucket policy
type DeleteBucketPolicyRequest struct {
	BucketID string `json:"bucket_id"`
	Name     string `json:"name"`
}

// DeleteBucketPolicy deletes a policy for a storage bucket
func (sc *StorageController) DeleteBucketPolicy(c *gin.Context) {
	var req DeleteBucketPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BucketID == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	sql := fmt.Sprintf(`
		DELETE FROM storage.policies
		WHERE 
			bucket_id = '%s' AND 
			name = '%s'
	`, req.BucketID, req.Name)

	var result interface{}
	err := sc.supabase.Functions.Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Policy '%s' for bucket '%s' deleted successfully", req.Name, req.BucketID),
	})
}
