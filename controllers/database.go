package controllers

import (
	"fmt"
	"net/http"

	"github.com/dirgocs/supabase-self-hosted-mcp/supabase"
	"github.com/dirgocs/supabase-self-hosted-mcp/types"
	"github.com/dirgocs/supabase-self-hosted-mcp/utils"
	"github.com/gin-gonic/gin"
)



// DatabaseController handles database-related operations
type DatabaseController struct {
	supabase *supabase.SupabaseClientExtended
}

// NewDatabaseController creates a new database controller
func NewDatabaseController(client *supabase.SupabaseClientExtended) *DatabaseController {
	return &DatabaseController{
		supabase: client,
	}
}

// ExecuteQueryRequest represents the request body for executing a query
type ExecuteQueryRequest struct {
	Query string `json:"query"`
}

// ExecuteQuery executes a SQL query (read-only for security)
func (dc *DatabaseController) ExecuteQuery(c *gin.Context) {
	var req ExecuteQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	// Check if the query is read-only
	if !utils.IsReadOnlyQuery(req.Query) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only read-only queries are allowed through this endpoint for security reasons",
		})
		return
	}

	// Execute query through Functions API
	var result interface{}
	// Use the Functions API to execute the query
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": req.Query,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Unable to execute query. You may need to create a custom function 'execute_sql' in your Supabase instance.",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDatabaseSchemaRequest represents the request body for getting database schema
type GetDatabaseSchemaRequest struct {
	Schema string `json:"schema"`
}

// SchemaData represents the processed schema data structure
type SchemaData map[string]map[string]TableSchema

// TableSchema represents the schema of a table
type TableSchema struct {
	Columns     []DatabaseColumn `json:"columns"`
	PrimaryKeys []string        `json:"primary_keys"`
	ForeignKeys []ForeignKey    `json:"foreign_keys"`
}

// DatabaseColumn represents a column in a table schema
type DatabaseColumn struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	NotNull      bool   `json:"not_null"`
	DefaultValue string `json:"default_value"`
}

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	Column     string            `json:"column"`
	References ForeignKeyReference `json:"references"`
}

// ForeignKeyReference represents the reference part of a foreign key
type ForeignKeyReference struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Column string `json:"column"`
}

// SchemaItem represents a raw schema item from the database
type SchemaItem struct {
	SchemaName      string `json:"schema_name"`
	TableName       string `json:"table_name"`
	ColumnName      string `json:"column_name"`
	DataType        string `json:"data_type"`
	NotNull         bool   `json:"not_null"`
	DefaultValue    string `json:"default_value"`
	IsPrimaryKey    bool   `json:"is_primary_key"`
	IsForeignKey    bool   `json:"is_foreign_key"`
	ReferenceSchema string `json:"reference_schema"`
	ReferenceTable  string `json:"reference_table"`
	ReferenceColumn string `json:"reference_column"`
}

// GetDatabaseSchema gets database schema information
func (dc *DatabaseController) GetDatabaseSchema(c *gin.Context) {
	var req GetDatabaseSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		SELECT 
			n.nspname AS schema_name,
			c.relname AS table_name,
			a.attname AS column_name,
			format_type(a.atttypid, a.atttypmod) AS data_type,
			a.attnotnull AS not_null,
			pg_get_expr(d.adbin, d.adrelid) AS default_value,
			CASE 
				WHEN co.contype = 'p' THEN true
				ELSE false
			END AS is_primary_key,
			CASE 
				WHEN co.contype = 'f' THEN true
				ELSE false
			END AS is_foreign_key,
			CASE 
				WHEN co.contype = 'f' THEN obj_schema.nspname
				ELSE NULL
			END AS reference_schema,
			CASE 
				WHEN co.contype = 'f' THEN ref_class.relname
				ELSE NULL
			END AS reference_table,
			CASE 
				WHEN co.contype = 'f' THEN ref_attr.attname
				ELSE NULL
			END AS reference_column
		FROM pg_attribute a
		JOIN pg_class c ON a.attrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		LEFT JOIN pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
		LEFT JOIN pg_constraint co ON (
			co.conrelid = c.oid AND 
			a.attnum = ANY(co.conkey) AND 
			(co.contype = 'p' OR co.contype = 'f')
		)
		LEFT JOIN pg_class ref_class ON co.confrelid = ref_class.oid
		LEFT JOIN pg_namespace obj_schema ON ref_class.relnamespace = obj_schema.oid
		LEFT JOIN pg_attribute ref_attr ON (
			ref_attr.attrelid = co.confrelid AND 
			ref_attr.attnum = co.confkey[array_position(co.conkey, a.attnum)]
		)
		WHERE a.attnum > 0 AND NOT a.attisdropped AND c.relkind = 'r'
	`

	if req.Schema != "" {
		query += fmt.Sprintf(" AND n.nspname = '%s'", req.Schema)
	} else {
		query += " AND n.nspname NOT IN ('pg_catalog', 'information_schema')"
	}

	query += " ORDER BY n.nspname, c.relname, a.attnum"

	// Execute query through RPC
	var rawData []types.SchemaItem
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": query,
	}, &rawData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process data into a more user-friendly structure
	schemaData := make(types.SchemaData)

	for _, item := range rawData {
		schemaName := item.SchemaName
		tableName := item.TableName

		if _, exists := schemaData[schemaName]; !exists {
			schemaData[schemaName] = make(map[string]types.TableSchema)
		}

		if _, exists := schemaData[schemaName][tableName]; !exists {
			schemaData[schemaName][tableName] = types.TableSchema{
				Columns:     []types.DatabaseColumn{},
				PrimaryKeys: []string{},
				ForeignKeys: []types.ForeignKey{},
			}
		}

		// Add column
		column := types.DatabaseColumn{
			Name:         item.ColumnName,
			Type:         item.DataType,
			NotNull:      item.NotNull,
			DefaultValue: item.DefaultValue,
		}

		tableSchema := schemaData[schemaName][tableName]
		tableSchema.Columns = append(tableSchema.Columns, column)

		// Add primary key
		if item.IsPrimaryKey {
			tableSchema.PrimaryKeys = append(tableSchema.PrimaryKeys, item.ColumnName)
		}

		// Add foreign key
		if item.IsForeignKey {
			tableSchema.ForeignKeys = append(tableSchema.ForeignKeys, types.ForeignKey{
				Column: item.ColumnName,
				References: types.ForeignKeyReference{
					Table: item.ReferencedTable,
					Column: item.ReferencedColumn,
				},
			})
		}

		schemaData[schemaName][tableName] = tableSchema
	}

	c.JSON(http.StatusOK, schemaData)
}

// CreateSchemaRequest represents the request body for creating a schema
type CreateSchemaRequest struct {
	Name string `json:"name"`
}

// CreateSchema creates a new schema
func (dc *DatabaseController) CreateSchema(c *gin.Context) {
	var req CreateSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema name is required"})
		return
	}

	sql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, req.Name)

	var result interface{}
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Schema '%s' created successfully", req.Name),
	})
}

// DeleteSchemaRequest represents the request body for deleting a schema
type DeleteSchemaRequest struct {
	Name    string `json:"name"`
	Cascade bool   `json:"cascade"`
}

// DeleteSchema deletes a schema
func (dc *DatabaseController) DeleteSchema(c *gin.Context) {
	var req DeleteSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema name is required"})
		return
	}

	sql := fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s"`, req.Name)

	if req.Cascade {
		sql += " CASCADE"
	}

	var result interface{}
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Schema '%s' deleted successfully", req.Name),
	})
}

// GetRLSPoliciesRequest represents the request body for getting RLS policies
type GetRLSPoliciesRequest struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
}

// RLSPolicy represents an RLS policy
type RLSPolicy struct {
	SchemaName       string   `json:"schema_name"`
	TableName        string   `json:"table_name"`
	PolicyName       string   `json:"policy_name"`
	PolicyType       string   `json:"policy_type"`
	Command          string   `json:"command"`
	Expression       string   `json:"expression"`
	CheckExpression  string   `json:"check_expression"`
	Roles            []string `json:"roles"`
}

// GetRLSPolicies gets RLS policies
func (dc *DatabaseController) GetRLSPolicies(c *gin.Context) {
	var req GetRLSPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schema := req.Schema
	if schema == "" {
		schema = "public"
	}

	query := `
		SELECT 
			n.nspname AS schema_name,
			c.relname AS table_name,
			p.polname AS policy_name,
			CASE 
				WHEN p.polpermissive THEN 'PERMISSIVE'
				ELSE 'RESTRICTIVE'
			END AS policy_type,
			CASE p.polcmd
				WHEN 'r' THEN 'SELECT'
				WHEN 'a' THEN 'INSERT'
				WHEN 'w' THEN 'UPDATE'
				WHEN 'd' THEN 'DELETE'
				WHEN '*' THEN 'ALL'
			END AS command,
			pg_get_expr(p.polqual, p.polrelid) AS expression,
			pg_get_expr(p.polwithcheck, p.polrelid) AS check_expression,
			ARRAY(
				SELECT rolname 
				FROM pg_roles 
				WHERE pg_has_role(oid, p.polroles, 'MEMBER')
			) AS roles
		FROM pg_policy p
		JOIN pg_class c ON p.polrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = $1
	`

	params := []interface{}{schema}

	if req.Table != "" {
		query += " AND c.relname = $2"
		params = append(params, req.Table)
	}

	query += " ORDER BY n.nspname, c.relname, p.polname"

	var policies []types.RLSPolicy

	// Execute the query to get RLS policies
	err := dc.supabase.Functions().Invoke("get_rls_policies", map[string]interface{}{
		"query": query,
		"params": params,
	}, &policies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Unable to execute query. You may need to create a custom function 'execute_sql' in your Supabase instance.",
		})
		return
	}

	c.JSON(http.StatusOK, policies)
}

// CreateRLSPolicyRequest represents the request body for creating an RLS policy
type CreateRLSPolicyRequest struct {
	Schema     string `json:"schema"`
	Table      string `json:"table"`
	Name       string `json:"name"`
	Operation  string `json:"operation"`
	Definition string `json:"definition"`
	Check      string `json:"check"`
	Role       string `json:"role"`
}

// CreateRLSPolicy creates a new RLS policy
func (dc *DatabaseController) CreateRLSPolicy(c *gin.Context) {
	var req CreateRLSPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Table == "" || req.Name == "" || req.Operation == "" || req.Definition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	schema := req.Schema
	if schema == "" {
		schema = "public"
	}

	role := req.Role
	if role == "" {
		role = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, schema, req.Table)

	sql := fmt.Sprintf(`CREATE POLICY "%s" ON %s 
               FOR %s 
               TO "%s" 
               USING (%s)`, req.Name, tableIdentifier, req.Operation, role, req.Definition)

	// Add WITH CHECK clause for operations that need it
	if req.Check != "" && (req.Operation == "INSERT" || req.Operation == "UPDATE" || req.Operation == "ALL") {
		sql += fmt.Sprintf(` WITH CHECK (%s)`, req.Check)
	}

	var result interface{}
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("RLS policy '%s' created on %s", req.Name, tableIdentifier),
	})
}

// UpdateRLSPolicyRequest represents the request body for updating an RLS policy
type UpdateRLSPolicyRequest struct {
	Schema     string `json:"schema"`
	Table      string `json:"table"`
	Name       string `json:"name"`
	Operation  string `json:"operation"`
	Definition string `json:"definition"`
	Check      string `json:"check"`
}

// UpdateRLSPolicy updates an existing RLS policy
func (dc *DatabaseController) UpdateRLSPolicy(c *gin.Context) {
	var req UpdateRLSPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Table == "" || req.Name == "" || req.Definition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	schema := req.Schema
	if schema == "" {
		schema = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, schema, req.Table)

	// First drop the existing policy
	dropSQL := fmt.Sprintf(`DROP POLICY IF EXISTS "%s" ON %s`, req.Name, tableIdentifier)
	var dropResult interface{}
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": dropSQL,
	}, &dropResult)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Then recreate it with new parameters
	operation := req.Operation
	if operation == "" {
		operation = "ALL"
	}

	createSQL := fmt.Sprintf(`CREATE POLICY "%s" ON %s 
                      FOR %s 
                      USING (%s)`, req.Name, tableIdentifier, operation, req.Definition)

	// Add WITH CHECK clause for operations that need it
	if req.Check != "" && (operation == "INSERT" || operation == "UPDATE" || operation == "ALL") {
		createSQL += fmt.Sprintf(` WITH CHECK (%s)`, req.Check)
	}

	var createResult interface{}
	err = dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": createSQL,
	}, &createResult)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("RLS policy '%s' updated on %s", req.Name, tableIdentifier),
	})
}

// DeleteRLSPolicyRequest represents the request body for deleting an RLS policy
type DeleteRLSPolicyRequest struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Name   string `json:"name"`
}

// DeleteRLSPolicy deletes an RLS policy
func (dc *DatabaseController) DeleteRLSPolicy(c *gin.Context) {
	var req DeleteRLSPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Table == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	schema := req.Schema
	if schema == "" {
		schema = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, schema, req.Table)
	sql := fmt.Sprintf(`DROP POLICY IF EXISTS "%s" ON %s`, req.Name, tableIdentifier)

	var result interface{}
	err := dc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("RLS policy '%s' deleted from %s", req.Name, tableIdentifier),
	})
}
