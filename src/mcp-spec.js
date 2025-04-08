// MCP Server specification
const mcpSpec = {
  functions: [
    {
      name: 'query_table',
      description: 'Query a specific table with schema selection and where clause support',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          },
          table: {
            type: 'string',
            description: 'Name of the table to query'
          },
          select: {
            type: 'string',
            description: 'Comma-separated list of columns to select (optional, defaults to *)'
          },
          where: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                column: {
                  type: 'string',
                  description: 'Column name'
                },
                operator: {
                  type: 'string',
                  enum: ['eq', 'neq', 'gt', 'gte', 'lt', 'lte', 'like', 'ilike', 'is'],
                  description: 'Comparison operator'
                },
                value: {
                  type: 'any',
                  description: 'Value to compare against'
                }
              },
              required: ['column', 'operator', 'value']
            },
            description: 'Array of where conditions (optional)'
          }
        },
        required: ['table']
      }
    },
    {
      name: 'generate_types',
      description: 'Generate TypeScript types for your Supabase database schema',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          }
        }
      }
    },
    {
      name: 'list_tables',
      description: 'List all tables in a specific schema',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          }
        }
      }
    },
    {
      name: 'execute_query',
      description: 'Execute a raw SQL query (with security restrictions)',
      parameters: {
        type: 'object',
        properties: {
          query: {
            type: 'string',
            description: 'SQL query to execute (read-only operations only)'
          }
        },
        required: ['query']
      }
    },
    // RLS Policies
    {
      name: 'get_rls_policies',
      description: 'Get RLS policies for a table or all tables',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          },
          table: {
            type: 'string',
            description: 'Table name (optional, if not provided returns policies for all tables)'
          }
        }
      }
    },
    {
      name: 'create_rls_policy',
      description: 'Create a new RLS policy',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          },
          table: {
            type: 'string',
            description: 'Table name'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          },
          operation: {
            type: 'string',
            enum: ['SELECT', 'INSERT', 'UPDATE', 'DELETE', 'ALL'],
            description: 'Operation type that the policy applies to'
          },
          definition: {
            type: 'string',
            description: 'Policy definition (using expression syntax)'
          },
          check: {
            type: 'string',
            description: 'Optional check expression for INSERT/UPDATE operations'
          },
          role: {
            type: 'string',
            description: 'Optional role name (defaults to public)'
          }
        },
        required: ['table', 'name', 'operation', 'definition']
      }
    },
    {
      name: 'update_rls_policy',
      description: 'Update an existing RLS policy',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          },
          table: {
            type: 'string',
            description: 'Table name'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          },
          operation: {
            type: 'string',
            enum: ['SELECT', 'INSERT', 'UPDATE', 'DELETE', 'ALL'],
            description: 'Operation type that the policy applies to'
          },
          definition: {
            type: 'string',
            description: 'Policy definition (using expression syntax)'
          },
          check: {
            type: 'string',
            description: 'Optional check expression for INSERT/UPDATE operations'
          }
        },
        required: ['table', 'name', 'definition']
      }
    },
    {
      name: 'delete_rls_policy',
      description: 'Delete an RLS policy',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Database schema (optional, defaults to public)'
          },
          table: {
            type: 'string',
            description: 'Table name'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          }
        },
        required: ['table', 'name']
      }
    },
    
    // Edge Functions
    {
      name: 'get_edge_functions',
      description: 'Get all edge functions or a specific one',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Function name (optional, if not provided returns all functions)'
          }
        }
      }
    },
    {
      name: 'create_edge_function',
      description: 'Create a new edge function',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Function name'
          },
          code: {
            type: 'string',
            description: 'Function code (JavaScript/TypeScript)'
          },
          verify_jwt: {
            type: 'boolean',
            description: 'Whether to verify JWT (optional, defaults to false)'
          },
          import_map: {
            type: 'object',
            description: 'Optional import map for the function'
          }
        },
        required: ['name', 'code']
      }
    },
    {
      name: 'update_edge_function',
      description: 'Update an existing edge function',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Function name'
          },
          code: {
            type: 'string',
            description: 'Function code (JavaScript/TypeScript)'
          },
          verify_jwt: {
            type: 'boolean',
            description: 'Whether to verify JWT (optional)'
          },
          import_map: {
            type: 'object',
            description: 'Optional import map for the function'
          }
        },
        required: ['name', 'code']
      }
    },
    {
      name: 'delete_edge_function',
      description: 'Delete an edge function',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Function name'
          }
        },
        required: ['name']
      }
    },
    {
      name: 'deploy_edge_function',
      description: 'Deploy an edge function',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Function name'
          }
        },
        required: ['name']
      }
    },
    
    // Database Schema
    {
      name: 'get_database_schema',
      description: 'Get database schema',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Schema name (optional, defaults to all schemas)'
          }
        }
      }
    },
    {
      name: 'create_schema',
      description: 'Create a new schema',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Schema name'
          }
        },
        required: ['name']
      }
    },
    {
      name: 'delete_schema',
      description: 'Delete a schema',
      parameters: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: 'Schema name'
          },
          cascade: {
            type: 'boolean',
            description: 'Whether to cascade the deletion (optional, defaults to false)'
          }
        },
        required: ['name']
      }
    },
    
    // Tables
    {
      name: 'create_table',
      description: 'Create a new table',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Schema name (optional, defaults to public)'
          },
          name: {
            type: 'string',
            description: 'Table name'
          },
          columns: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                name: {
                  type: 'string',
                  description: 'Column name'
                },
                type: {
                  type: 'string',
                  description: 'Column data type'
                },
                nullable: {
                  type: 'boolean',
                  description: 'Whether the column can be null (optional, defaults to true)'
                },
                default_value: {
                  type: 'string',
                  description: 'Default value (optional)'
                },
                primary_key: {
                  type: 'boolean',
                  description: 'Whether the column is a primary key (optional, defaults to false)'
                },
                unique: {
                  type: 'boolean',
                  description: 'Whether the column value must be unique (optional, defaults to false)'
                },
                references: {
                  type: 'object',
                  properties: {
                    table: {
                      type: 'string',
                      description: 'Referenced table'
                    },
                    column: {
                      type: 'string',
                      description: 'Referenced column'
                    }
                  },
                  description: 'Foreign key reference (optional)'
                }
              },
              required: ['name', 'type']
            },
            description: 'Table columns'
          },
          enable_rls: {
            type: 'boolean',
            description: 'Whether to enable RLS on the table (optional, defaults to false)'
          }
        },
        required: ['name', 'columns']
      }
    },
    {
      name: 'alter_table',
      description: 'Alter a table (add/drop columns, rename)',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Schema name (optional, defaults to public)'
          },
          name: {
            type: 'string',
            description: 'Table name'
          },
          new_name: {
            type: 'string',
            description: 'New table name (optional)'
          },
          add_columns: {
            type: 'array',
            items: {
              type: 'object',
              properties: {
                name: {
                  type: 'string',
                  description: 'Column name'
                },
                type: {
                  type: 'string',
                  description: 'Column data type'
                },
                nullable: {
                  type: 'boolean',
                  description: 'Whether the column can be null (optional, defaults to true)'
                },
                default_value: {
                  type: 'string',
                  description: 'Default value (optional)'
                }
              },
              required: ['name', 'type']
            },
            description: 'Columns to add (optional)'
          },
          drop_columns: {
            type: 'array',
            items: {
              type: 'string',
              description: 'Column name to drop'
            },
            description: 'Columns to drop (optional)'
          },
          enable_rls: {
            type: 'boolean',
            description: 'Whether to enable RLS on the table (optional)'
          }
        },
        required: ['name']
      }
    },
    {
      name: 'drop_table',
      description: 'Drop a table',
      parameters: {
        type: 'object',
        properties: {
          schema: {
            type: 'string',
            description: 'Schema name (optional, defaults to public)'
          },
          name: {
            type: 'string',
            description: 'Table name'
          },
          cascade: {
            type: 'boolean',
            description: 'Whether to cascade the deletion (optional, defaults to false)'
          }
        },
        required: ['name']
      }
    },
    
    // Storage Buckets
    {
      name: 'get_buckets',
      description: 'Get all storage buckets or a specific one',
      parameters: {
        type: 'object',
        properties: {
          id: {
            type: 'string',
            description: 'Bucket ID (optional, if not provided returns all buckets)'
          }
        }
      }
    },
    {
      name: 'create_bucket',
      description: 'Create a new storage bucket',
      parameters: {
        type: 'object',
        properties: {
          id: {
            type: 'string',
            description: 'Bucket ID'
          },
          name: {
            type: 'string',
            description: 'Bucket name (optional, defaults to ID)'
          },
          public: {
            type: 'boolean',
            description: 'Whether the bucket is public (optional, defaults to false)'
          },
          file_size_limit: {
            type: 'number',
            description: 'File size limit in bytes (optional)'
          },
          allowed_mime_types: {
            type: 'array',
            items: {
              type: 'string'
            },
            description: 'Allowed MIME types (optional)'
          }
        },
        required: ['id']
      }
    },
    {
      name: 'update_bucket',
      description: 'Update a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          id: {
            type: 'string',
            description: 'Bucket ID'
          },
          public: {
            type: 'boolean',
            description: 'Whether the bucket is public (optional)'
          },
          file_size_limit: {
            type: 'number',
            description: 'File size limit in bytes (optional)'
          },
          allowed_mime_types: {
            type: 'array',
            items: {
              type: 'string'
            },
            description: 'Allowed MIME types (optional)'
          }
        },
        required: ['id']
      }
    },
    {
      name: 'delete_bucket',
      description: 'Delete a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          id: {
            type: 'string',
            description: 'Bucket ID'
          }
        },
        required: ['id']
      }
    },
    
    // Bucket Policies
    {
      name: 'get_bucket_policies',
      description: 'Get policies for a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          bucket_id: {
            type: 'string',
            description: 'Bucket ID'
          }
        },
        required: ['bucket_id']
      }
    },
    {
      name: 'create_bucket_policy',
      description: 'Create a new policy for a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          bucket_id: {
            type: 'string',
            description: 'Bucket ID'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          },
          operation: {
            type: 'string',
            enum: ['SELECT', 'INSERT', 'UPDATE', 'DELETE'],
            description: 'Operation type'
          },
          definition: {
            type: 'string',
            description: 'Policy definition (using expression syntax)'
          },
          role: {
            type: 'string',
            description: 'Optional role name (defaults to public)'
          }
        },
        required: ['bucket_id', 'name', 'operation', 'definition']
      }
    },
    {
      name: 'update_bucket_policy',
      description: 'Update a policy for a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          bucket_id: {
            type: 'string',
            description: 'Bucket ID'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          },
          definition: {
            type: 'string',
            description: 'Policy definition (using expression syntax)'
          }
        },
        required: ['bucket_id', 'name', 'definition']
      }
    },
    {
      name: 'delete_bucket_policy',
      description: 'Delete a policy for a storage bucket',
      parameters: {
        type: 'object',
        properties: {
          bucket_id: {
            type: 'string',
            description: 'Bucket ID'
          },
          name: {
            type: 'string',
            description: 'Policy name'
          }
        },
        required: ['bucket_id', 'name']
      }
    }
  ]
};

module.exports = mcpSpec;
