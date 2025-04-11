package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dirgocs/supabase-self-hosted-mcp/supabase"
	"github.com/dirgocs/supabase-self-hosted-mcp/utils"
	"github.com/gin-gonic/gin"
)

// TableController handles table-related operations
type TableController struct {
	supabase *supabase.SupabaseClientExtended
}

// NewTableController creates a new table controller
func NewTableController(client *supabase.SupabaseClientExtended) *TableController {
	return &TableController{
		supabase: client,
	}
}

// WhereCondition represents a where condition for table queries
type WhereCondition struct {
	Column   string      `json:"column"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// QueryTableRequest represents the request body for querying a table
type QueryTableRequest struct {
	Schema string           `json:"schema"`
	Table  string           `json:"table"`
	Select string           `json:"select"`
	Where  []WhereCondition `json:"where"`
}

// QueryTable queries a specific table with filters
func (tc *TableController) QueryTable(c *gin.Context) {
	var req QueryTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	// Set defaults
	if req.Schema == "" {
		req.Schema = "public"
	}
	if req.Select == "" {
		req.Select = "*"
	}

	// Build the query
	query := tc.supabase.DB.From(req.Table).Select(req.Select)

	// Apply where conditions
	for _, condition := range req.Where {
		switch condition.Operator {
		case "eq":
			if strValue, ok := condition.Value.(string); ok {
				query.Eq(condition.Column, strValue)
			} else {
				query.Eq(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "neq":
			if strValue, ok := condition.Value.(string); ok {
				query.Neq(condition.Column, strValue)
			} else {
				query.Neq(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "gt":
			if strValue, ok := condition.Value.(string); ok {
				query.Gt(condition.Column, strValue)
			} else {
				query.Gt(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "gte":
			if strValue, ok := condition.Value.(string); ok {
				query.Gte(condition.Column, strValue)
			} else {
				query.Gte(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "lt":
			if strValue, ok := condition.Value.(string); ok {
				query.Lt(condition.Column, strValue)
			} else {
				query.Lt(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "lte":
			if strValue, ok := condition.Value.(string); ok {
				query.Lte(condition.Column, strValue)
			} else {
				query.Lte(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "like":
			if strValue, ok := condition.Value.(string); ok {
				query.Like(condition.Column, strValue)
			} else {
				query.Like(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "ilike":
			if strValue, ok := condition.Value.(string); ok {
				query.Ilike(condition.Column, strValue)
			} else {
				query.Ilike(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		case "is":
			// Handle IS operator with proper type conversion
			if condition.Value == nil {
				// For IS NULL, use an empty string as a workaround
				query.Is(condition.Column, "")
			} else if strValue, ok := condition.Value.(string); ok {
				query.Is(condition.Column, strValue)
			} else {
				query.Is(condition.Column, fmt.Sprintf("%v", condition.Value))
			}
		default:
			// Ignore unknown operators
		}
	}

	var result []map[string]interface{}
	err := query.Execute(&result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateTypesRequest represents the request body for generating types
type GenerateTypesRequest struct {
	Schema string `json:"schema"`
}

// SchemaColumn represents a column in the database schema
type SchemaColumn struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable string `json:"is_nullable"`
}

// SchemaTable represents a table in the database schema
type SchemaTable struct {
	TableName string         `json:"table_name"`
	Columns   []SchemaColumn `json:"columns"`
}

// GenerateTypes generates TypeScript types for a schema
func (tc *TableController) GenerateTypes(c *gin.Context) {
	var req GenerateTypesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default schema
	if req.Schema == "" {
		req.Schema = "public"
	}

	// Try to get schema information using RPC
	var tables []SchemaTable
	err := tc.supabase.Functions().Invoke("get_schema_information", map[string]interface{}{
		"p_schema": req.Schema,
	}, &tables)

	if err != nil {
		// If RPC fails, try direct query
		var schemaData []map[string]interface{}
		// Use the DB object to execute a SQL query
		query := fmt.Sprintf(`SELECT tablename FROM pg_tables WHERE schemaname = '%s'`, req.Schema)
		err = tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
			"query": query,
		}, &schemaData)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Unable to generate types. You may need to create a custom function to retrieve schema information in your Supabase instance.",
			})
			return
		}

		// Build simple types based on tables found
		var typesOutput strings.Builder
		typesOutput.WriteString(fmt.Sprintf("// TypeScript types for schema: %s\n\n", req.Schema))

		for _, tableInfo := range schemaData {
			tableName := tableInfo["tablename"].(string)

			typesOutput.WriteString(fmt.Sprintf("export interface %s {\n", utils.ToPascalCase(tableName)))
			typesOutput.WriteString("  // Add your column definitions here\n")
			typesOutput.WriteString("  id: string; // Assuming primary key\n")
			typesOutput.WriteString("  created_at?: string; // Common timestamp field\n")
			typesOutput.WriteString("}\n\n")
		}

		c.JSON(http.StatusOK, gin.H{"types": typesOutput.String()})
		return
	}

	// If RPC was successful, build detailed types
	var typesOutput strings.Builder
	typesOutput.WriteString(fmt.Sprintf("// TypeScript types for schema: %s\n\n", req.Schema))

	for _, tableInfo := range tables {
		tableName := tableInfo.TableName

		typesOutput.WriteString(fmt.Sprintf("export interface %s {\n", utils.ToPascalCase(tableName)))

		for _, column := range tableInfo.Columns {
			columnName := column.ColumnName
			var columnType string

			// Map PostgreSQL types to TypeScript types
			switch strings.ToLower(column.DataType) {
			case "integer", "numeric", "decimal", "real", "double precision", "smallint", "bigint":
				columnType = "number"
			case "text", "character varying", "character", "varchar", "char", "uuid", "date", "time", "timestamp", "timestamptz":
				columnType = "string"
			case "boolean":
				columnType = "boolean"
			case "jsonb", "json":
				columnType = "Record<string, any>"
			case "array":
				columnType = "any[]"
			default:
				columnType = "any"
			}

			isNullable := column.IsNullable == "YES"
			nullableModifier := ""
			if isNullable {
				nullableModifier = "?"
			}

			typesOutput.WriteString(fmt.Sprintf("  %s%s: %s;\n", columnName, nullableModifier, columnType))
		}

		typesOutput.WriteString("}\n\n")
	}

	c.JSON(http.StatusOK, gin.H{"types": typesOutput.String()})
}

// ListTablesRequest represents the request body for listing tables
type ListTablesRequest struct {
	Schema string `json:"schema"`
}

// ListTables lists tables in a schema
func (tc *TableController) ListTables(c *gin.Context) {
	var req ListTablesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default schema
	if req.Schema == "" {
		req.Schema = "public"
	}

	// Try using RPC first
	var tables []map[string]interface{}
	err := tc.supabase.Functions().Invoke("list_tables", map[string]interface{}{
		"p_schema": req.Schema,
	}, &tables)

	if err == nil && tables != nil {
		c.JSON(http.StatusOK, tables)
		return
	}

	// Alternative: direct query
	query := fmt.Sprintf(`
		SELECT tablename AS table_name 
		FROM pg_tables 
		WHERE schemaname = '%s'
	`, req.Schema)

	var queryResult []map[string]interface{}
	err = tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": query,
	}, &queryResult)

	if err != nil {
		// If RPC also fails, try checking common tables
		var commonTables []map[string]string
		commonTableNames := []string{"profiles", "products", "users", "categories", "orders", "items"}

		for _, tableName := range commonTableNames {
			// Use a simple query to check if the table exists
			query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
			var result []map[string]interface{}
			err := tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
				"query": query,
			}, &result)

			if err == nil {
				commonTables = append(commonTables, map[string]string{"table_name": tableName})
			}
		}

		c.JSON(http.StatusOK, commonTables)
		return
	}

	c.JSON(http.StatusOK, queryResult)
}

// Column represents a column definition for table creation
type Column struct {
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Nullable     *bool      `json:"nullable"`
	DefaultValue string     `json:"default_value"`
	PrimaryKey   bool       `json:"primary_key"`
	Unique       bool       `json:"unique"`
	References   *Reference `json:"references"`
}

// Reference represents a foreign key reference
type Reference struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// CreateTableRequest represents the request body for creating a table
type CreateTableRequest struct {
	Schema    string   `json:"schema"`
	Name      string   `json:"name"`
	Columns   []Column `json:"columns"`
	EnableRLS bool     `json:"enable_rls"`
}

// CreateTable creates a new table
func (tc *TableController) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || len(req.Columns) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name and columns are required"})
		return
	}

	// Set default schema
	if req.Schema == "" {
		req.Schema = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, req.Schema, req.Name)

	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableIdentifier))

	// Add columns
	var columnDefinitions []string
	for _, column := range req.Columns {
		def := fmt.Sprintf(`  "%s" %s`, column.Name, column.Type)

		// Not null
		if column.Nullable != nil && !*column.Nullable {
			def += " NOT NULL"
		}

		// Default value
		if column.DefaultValue != "" {
			def += fmt.Sprintf(" DEFAULT %s", column.DefaultValue)
		}

		// Primary key
		if column.PrimaryKey {
			def += " PRIMARY KEY"
		}

		// Unique
		if column.Unique {
			def += " UNIQUE"
		}

		// References (foreign key)
		if column.References != nil {
			def += fmt.Sprintf(` REFERENCES "%s" ("%s")`, column.References.Table, column.References.Column)
		}

		columnDefinitions = append(columnDefinitions, def)
	}

	sql.WriteString(strings.Join(columnDefinitions, ",\n"))
	sql.WriteString("\n)")

	// Enable RLS if requested
	if req.EnableRLS {
		sql.WriteString(fmt.Sprintf(";\nALTER TABLE %s ENABLE ROW LEVEL SECURITY", tableIdentifier))
	}

	var result interface{}
	err := tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql.String(),
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Table '%s' created successfully in schema '%s'", req.Name, req.Schema),
	})
}

// AlterTableRequest represents the request body for altering a table
type AlterTableRequest struct {
	Schema      string   `json:"schema"`
	Name        string   `json:"name"`
	NewName     string   `json:"new_name"`
	AddColumns  []Column `json:"add_columns"`
	DropColumns []string `json:"drop_columns"`
	EnableRLS   *bool    `json:"enable_rls"`
}

// OperationResult represents the result of a database operation
type OperationResult struct {
	Success bool   `json:"success"`
	Query   string `json:"query"`
	Error   string `json:"error,omitempty"`
}

// AlterTable alters an existing table
func (tc *TableController) AlterTable(c *gin.Context) {
	var req AlterTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	// Set default schema
	if req.Schema == "" {
		req.Schema = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, req.Schema, req.Name)
	var sqls []string

	// Rename table
	if req.NewName != "" {
		sqls = append(sqls, fmt.Sprintf(`ALTER TABLE %s RENAME TO "%s"`, tableIdentifier, req.NewName))
	}

	// Add columns
	if len(req.AddColumns) > 0 {
		for _, column := range req.AddColumns {
			sql := fmt.Sprintf(`ALTER TABLE %s ADD COLUMN "%s" %s`, tableIdentifier, column.Name, column.Type)

			if column.Nullable != nil && !*column.Nullable {
				sql += " NOT NULL"
			}

			if column.DefaultValue != "" {
				sql += fmt.Sprintf(" DEFAULT %s", column.DefaultValue)
			}

			sqls = append(sqls, sql)
		}
	}

	// Drop columns
	if len(req.DropColumns) > 0 {
		for _, columnName := range req.DropColumns {
			sqls = append(sqls, fmt.Sprintf(`ALTER TABLE %s DROP COLUMN "%s"`, tableIdentifier, columnName))
		}
	}

	// Enable/disable RLS
	if req.EnableRLS != nil {
		if *req.EnableRLS {
			sqls = append(sqls, fmt.Sprintf(`ALTER TABLE %s ENABLE ROW LEVEL SECURITY`, tableIdentifier))
		} else {
			sqls = append(sqls, fmt.Sprintf(`ALTER TABLE %s DISABLE ROW LEVEL SECURITY`, tableIdentifier))
		}
	}

	// Execute all queries
	var results []OperationResult

	for _, sql := range sqls {
		var result interface{}
		err := tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
			"query": sql,
		}, &result)

		if err != nil {
			results = append(results, OperationResult{
				Success: false,
				Query:   sql,
				Error:   err.Error(),
			})
		} else {
			results = append(results, OperationResult{
				Success: true,
				Query:   sql,
			})
		}
	}

	// Check if all operations were successful
	allSuccess := true
	for _, result := range results {
		if !result.Success {
			allSuccess = false
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    allSuccess,
		"operations": results,
		"message": func() string {
			if allSuccess {
				return fmt.Sprintf("Table '%s' altered successfully", req.Name)
			}
			return fmt.Sprintf("Some operations failed while altering table '%s'", req.Name)
		}(),
	})
}

// DropTableRequest represents the request body for dropping a table
type DropTableRequest struct {
	Schema  string `json:"schema"`
	Name    string `json:"name"`
	Cascade bool   `json:"cascade"`
}

// DropTable drops a table
func (tc *TableController) DropTable(c *gin.Context) {
	var req DropTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	// Set default schema
	if req.Schema == "" {
		req.Schema = "public"
	}

	tableIdentifier := fmt.Sprintf(`"%s"."%s"`, req.Schema, req.Name)
	sql := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableIdentifier)

	if req.Cascade {
		sql += ` CASCADE`
	}

	var result interface{}
	err := tc.supabase.Functions().Invoke("execute_sql", map[string]interface{}{
		"query": sql,
	}, &result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Table '%s' dropped successfully from schema '%s'", req.Name, req.Schema),
	})
}
