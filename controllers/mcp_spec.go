package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MCPSpec represents the MCP server specification
var MCPSpec = gin.H{
	"functions": []gin.H{
		{
			"name":        "query_table",
			"description": "Query a specific table with schema selection and where clause support",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
					"table": gin.H{
						"type":        "string",
						"description": "Name of the table to query",
					},
					"select": gin.H{
						"type":        "string",
						"description": "Comma-separated list of columns to select (optional, defaults to *)",
					},
					"where": gin.H{
						"type": "array",
						"items": gin.H{
							"type": "object",
							"properties": gin.H{
								"column": gin.H{
									"type":        "string",
									"description": "Column name",
								},
								"operator": gin.H{
									"type":        "string",
									"enum":        []string{"eq", "neq", "gt", "gte", "lt", "lte", "like", "ilike", "is"},
									"description": "Comparison operator",
								},
								"value": gin.H{
									"type":        "any",
									"description": "Value to compare against",
								},
							},
							"required": []string{"column", "operator", "value"},
						},
						"description": "Array of where conditions (optional)",
					},
				},
				"required": []string{"table"},
			},
		},
		{
			"name":        "generate_types",
			"description": "Generate TypeScript types for your Supabase database schema",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
				},
			},
		},
		{
			"name":        "list_tables",
			"description": "List all tables in a specific schema",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
				},
			},
		},
		{
			"name":        "execute_query",
			"description": "Execute a raw SQL query (with security restrictions)",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"query": gin.H{
						"type":        "string",
						"description": "SQL query to execute (read-only operations only)",
					},
				},
				"required": []string{"query"},
			},
		},
		// RLS Policies
		{
			"name":        "get_rls_policies",
			"description": "Get RLS policies for a table or all tables",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
					"table": gin.H{
						"type":        "string",
						"description": "Table name (optional, if not provided returns policies for all tables)",
					},
				},
			},
		},
		{
			"name":        "create_rls_policy",
			"description": "Create a new RLS policy",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
					"table": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
					"operation": gin.H{
						"type":        "string",
						"enum":        []string{"SELECT", "INSERT", "UPDATE", "DELETE", "ALL"},
						"description": "Operation type that the policy applies to",
					},
					"definition": gin.H{
						"type":        "string",
						"description": "Policy definition (using expression syntax)",
					},
					"check": gin.H{
						"type":        "string",
						"description": "Optional check expression for INSERT/UPDATE operations",
					},
					"role": gin.H{
						"type":        "string",
						"description": "Optional role name (defaults to public)",
					},
				},
				"required": []string{"table", "name", "operation", "definition"},
			},
		},
		{
			"name":        "update_rls_policy",
			"description": "Update an existing RLS policy",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
					"table": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
					"operation": gin.H{
						"type":        "string",
						"enum":        []string{"SELECT", "INSERT", "UPDATE", "DELETE", "ALL"},
						"description": "Operation type that the policy applies to",
					},
					"definition": gin.H{
						"type":        "string",
						"description": "Policy definition (using expression syntax)",
					},
					"check": gin.H{
						"type":        "string",
						"description": "Optional check expression for INSERT/UPDATE operations",
					},
				},
				"required": []string{"table", "name", "definition"},
			},
		},
		{
			"name":        "delete_rls_policy",
			"description": "Delete an RLS policy",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Database schema (optional, defaults to public)",
					},
					"table": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
				},
				"required": []string{"table", "name"},
			},
		},
		// Edge Functions
		{
			"name":        "get_edge_functions",
			"description": "Get all edge functions or a specific one",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Function name (optional, if not provided returns all functions)",
					},
				},
			},
		},
		{
			"name":        "create_edge_function",
			"description": "Create a new edge function",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Function name",
					},
					"code": gin.H{
						"type":        "string",
						"description": "Function code (JavaScript/TypeScript)",
					},
					"verify_jwt": gin.H{
						"type":        "boolean",
						"description": "Whether to verify JWT (optional, defaults to false)",
					},
					"import_map": gin.H{
						"type":        "object",
						"description": "Optional import map for the function",
					},
				},
				"required": []string{"name", "code"},
			},
		},
		{
			"name":        "update_edge_function",
			"description": "Update an existing edge function",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Function name",
					},
					"code": gin.H{
						"type":        "string",
						"description": "Function code (JavaScript/TypeScript)",
					},
					"verify_jwt": gin.H{
						"type":        "boolean",
						"description": "Whether to verify JWT (optional)",
					},
					"import_map": gin.H{
						"type":        "object",
						"description": "Optional import map for the function",
					},
				},
				"required": []string{"name", "code"},
			},
		},
		{
			"name":        "delete_edge_function",
			"description": "Delete an edge function",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Function name",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "deploy_edge_function",
			"description": "Deploy an edge function",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Function name",
					},
				},
				"required": []string{"name"},
			},
		},
		// Database Schema
		{
			"name":        "get_database_schema",
			"description": "Get database schema",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Schema name (optional, defaults to all schemas)",
					},
				},
			},
		},
		{
			"name":        "create_schema",
			"description": "Create a new schema",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Schema name",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "delete_schema",
			"description": "Delete a schema",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"name": gin.H{
						"type":        "string",
						"description": "Schema name",
					},
					"cascade": gin.H{
						"type":        "boolean",
						"description": "Whether to cascade the deletion (optional, defaults to false)",
					},
				},
				"required": []string{"name"},
			},
		},
		// Tables
		{
			"name":        "create_table",
			"description": "Create a new table",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Schema name (optional, defaults to public)",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"columns": gin.H{
						"type": "array",
						"items": gin.H{
							"type": "object",
							"properties": gin.H{
								"name": gin.H{
									"type":        "string",
									"description": "Column name",
								},
								"type": gin.H{
									"type":        "string",
									"description": "Column data type",
								},
								"nullable": gin.H{
									"type":        "boolean",
									"description": "Whether the column can be null (optional, defaults to true)",
								},
								"default_value": gin.H{
									"type":        "string",
									"description": "Default value (optional)",
								},
								"primary_key": gin.H{
									"type":        "boolean",
									"description": "Whether the column is a primary key (optional, defaults to false)",
								},
								"unique": gin.H{
									"type":        "boolean",
									"description": "Whether the column value must be unique (optional, defaults to false)",
								},
								"references": gin.H{
									"type": "object",
									"properties": gin.H{
										"table": gin.H{
											"type":        "string",
											"description": "Referenced table",
										},
										"column": gin.H{
											"type":        "string",
											"description": "Referenced column",
										},
									},
									"description": "Foreign key reference (optional)",
								},
							},
							"required": []string{"name", "type"},
						},
						"description": "Table columns",
					},
					"enable_rls": gin.H{
						"type":        "boolean",
						"description": "Whether to enable RLS on the table (optional, defaults to false)",
					},
				},
				"required": []string{"name", "columns"},
			},
		},
		{
			"name":        "alter_table",
			"description": "Alter a table (add/drop columns, rename)",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Schema name (optional, defaults to public)",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"new_name": gin.H{
						"type":        "string",
						"description": "New table name (optional)",
					},
					"add_columns": gin.H{
						"type": "array",
						"items": gin.H{
							"type": "object",
							"properties": gin.H{
								"name": gin.H{
									"type":        "string",
									"description": "Column name",
								},
								"type": gin.H{
									"type":        "string",
									"description": "Column data type",
								},
								"nullable": gin.H{
									"type":        "boolean",
									"description": "Whether the column can be null (optional, defaults to true)",
								},
								"default_value": gin.H{
									"type":        "string",
									"description": "Default value (optional)",
								},
							},
							"required": []string{"name", "type"},
						},
						"description": "Columns to add (optional)",
					},
					"drop_columns": gin.H{
						"type": "array",
						"items": gin.H{
							"type":        "string",
							"description": "Column name to drop",
						},
						"description": "Columns to drop (optional)",
					},
					"enable_rls": gin.H{
						"type":        "boolean",
						"description": "Whether to enable RLS on the table (optional)",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "drop_table",
			"description": "Drop a table",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"schema": gin.H{
						"type":        "string",
						"description": "Schema name (optional, defaults to public)",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Table name",
					},
					"cascade": gin.H{
						"type":        "boolean",
						"description": "Whether to cascade the deletion (optional, defaults to false)",
					},
				},
				"required": []string{"name"},
			},
		},
		// Storage Buckets
		{
			"name":        "get_buckets",
			"description": "Get all storage buckets or a specific one",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"id": gin.H{
						"type":        "string",
						"description": "Bucket ID (optional, if not provided returns all buckets)",
					},
				},
			},
		},
		{
			"name":        "create_bucket",
			"description": "Create a new storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Bucket name (optional, defaults to ID)",
					},
					"public": gin.H{
						"type":        "boolean",
						"description": "Whether the bucket is public (optional, defaults to false)",
					},
					"file_size_limit": gin.H{
						"type":        "number",
						"description": "File size limit in bytes (optional)",
					},
					"allowed_mime_types": gin.H{
						"type": "array",
						"items": gin.H{
							"type": "string",
						},
						"description": "Allowed MIME types (optional)",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "update_bucket",
			"description": "Update a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
					"public": gin.H{
						"type":        "boolean",
						"description": "Whether the bucket is public (optional)",
					},
					"file_size_limit": gin.H{
						"type":        "number",
						"description": "File size limit in bytes (optional)",
					},
					"allowed_mime_types": gin.H{
						"type": "array",
						"items": gin.H{
							"type": "string",
						},
						"description": "Allowed MIME types (optional)",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "delete_bucket",
			"description": "Delete a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
				},
				"required": []string{"id"},
			},
		},
		// Bucket Policies
		{
			"name":        "get_bucket_policies",
			"description": "Get policies for a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"bucket_id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
				},
				"required": []string{"bucket_id"},
			},
		},
		{
			"name":        "create_bucket_policy",
			"description": "Create a new policy for a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"bucket_id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
					"operation": gin.H{
						"type":        "string",
						"enum":        []string{"SELECT", "INSERT", "UPDATE", "DELETE"},
						"description": "Operation type",
					},
					"definition": gin.H{
						"type":        "string",
						"description": "Policy definition (using expression syntax)",
					},
					"role": gin.H{
						"type":        "string",
						"description": "Optional role name (defaults to public)",
					},
				},
				"required": []string{"bucket_id", "name", "operation", "definition"},
			},
		},
		{
			"name":        "update_bucket_policy",
			"description": "Update a policy for a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"bucket_id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
					"definition": gin.H{
						"type":        "string",
						"description": "Policy definition (using expression syntax)",
					},
				},
				"required": []string{"bucket_id", "name", "definition"},
			},
		},
		{
			"name":        "delete_bucket_policy",
			"description": "Delete a policy for a storage bucket",
			"parameters": gin.H{
				"type": "object",
				"properties": gin.H{
					"bucket_id": gin.H{
						"type":        "string",
						"description": "Bucket ID",
					},
					"name": gin.H{
						"type":        "string",
						"description": "Policy name",
					},
				},
				"required": []string{"bucket_id", "name"},
			},
		},
	},
}

// GetMCPSpecification returns the MCP specification
func GetMCPSpecification(c *gin.Context) {
	c.JSON(http.StatusOK, MCPSpec)
}
