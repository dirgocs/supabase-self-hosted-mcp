// Controllers for table-related operations
const { toPascalCase } = require('../utils');

/**
 * Query a specific table with filters
 */
async function queryTable(req, res, supabase) {
  try {
    const { schema = 'public', table, select = '*', where = [] } = req.body;
    
    if (!table) {
      return res.status(400).json({ error: 'Table name is required' });
    }

    let query = supabase
      .from(table)
      .select(select);

    // Apply where conditions
    where.forEach(condition => {
      const { column, operator, value } = condition;
      
      switch (operator) {
        case 'eq':
          query = query.eq(column, value);
          break;
        case 'neq':
          query = query.neq(column, value);
          break;
        case 'gt':
          query = query.gt(column, value);
          break;
        case 'gte':
          query = query.gte(column, value);
          break;
        case 'lt':
          query = query.lt(column, value);
          break;
        case 'lte':
          query = query.lte(column, value);
          break;
        case 'like':
          query = query.like(column, value);
          break;
        case 'ilike':
          query = query.ilike(column, value);
          break;
        case 'is':
          query = query.is(column, value);
          break;
        default:
          // Ignore unknown operators
          break;
      }
    });

    const { data, error } = await query;
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json(data);
  } catch (error) {
    console.error('Error querying table:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Generate TypeScript types for a schema
 */
async function generateTypes(req, res, supabase) {
  try {
    const { schema = 'public' } = req.body;
    
    // Try to get schema information using RPC
    const { data: tables, error: tablesError } = await supabase
      .rpc('get_schema_information', { p_schema: schema });
    
    if (tablesError) {
      // If RPC fails, try direct query
      const { data: schemaData, error: schemaError } = await supabase
        .from('pg_tables')
        .select('tablename')
        .eq('schemaname', schema);
      
      if (schemaError) {
        return res.status(400).json({ 
          error: schemaError.message, 
          message: "Unable to generate types. You may need to create a custom function to retrieve schema information in your Supabase instance."
        });
      }
      
      // Build simple types based on tables found
      let typesOutput = `// TypeScript types for schema: ${schema}\n\n`;
      
      for (const tableInfo of schemaData || []) {
        const tableName = tableInfo.tablename;
        
        typesOutput += `export interface ${toPascalCase(tableName)} {\n`;
        typesOutput += `  // Add your column definitions here\n`;
        typesOutput += `  id: string; // Assuming primary key\n`;
        typesOutput += `  created_at?: string; // Common timestamp field\n`;
        typesOutput += `}\n\n`;
      }
      
      return res.json({ types: typesOutput });
    }
    
    // If RPC was successful, build detailed types
    let typesOutput = `// TypeScript types for schema: ${schema}\n\n`;
    
    for (const tableInfo of tables || []) {
      const tableName = tableInfo.table_name;
      
      typesOutput += `export interface ${toPascalCase(tableName)} {\n`;
      
      for (const column of tableInfo.columns || []) {
        const columnName = column.column_name;
        let columnType = 'any';
        
        // Map PostgreSQL types to TypeScript types
        switch (column.data_type.toLowerCase()) {
          case 'integer':
          case 'numeric':
          case 'decimal':
          case 'real':
          case 'double precision':
          case 'smallint':
          case 'bigint':
            columnType = 'number';
            break;
          case 'text':
          case 'character varying':
          case 'character':
          case 'varchar':
          case 'char':
          case 'uuid':
          case 'date':
          case 'time':
          case 'timestamp':
          case 'timestamptz':
            columnType = 'string';
            break;
          case 'boolean':
            columnType = 'boolean';
            break;
          case 'jsonb':
          case 'json':
            columnType = 'Record<string, any>';
            break;
          case 'array':
            columnType = 'any[]';
            break;
          default:
            columnType = 'any';
            break;
        }
        
        const isNullable = column.is_nullable === 'YES';
        const nullableModifier = isNullable ? '?' : '';
        
        typesOutput += `  ${columnName}${nullableModifier}: ${columnType};\n`;
      }
      
      typesOutput += `}\n\n`;
    }
    
    return res.json({ types: typesOutput });
  } catch (error) {
    console.error('Error generating types:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * List tables in a schema
 */
async function listTables(req, res, supabase) {
  try {
    const { schema = 'public' } = req.body;
    
    // Try using RPC first
    try {
      const { data, error } = await supabase
        .rpc('list_tables', { p_schema: schema });
      
      if (!error && data) {
        return res.json(data);
      }
    } catch (rpcError) {
      console.log('RPC method not available, falling back to direct queries');
    }
    
    // Alternative: direct query
    const query = `
      SELECT tablename AS table_name 
      FROM pg_tables 
      WHERE schemaname = '${schema}'
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query });
    
    if (error) {
      // If RPC also fails, try checking common tables
      const tables = [];
      
      const commonTables = ['profiles', 'products', 'users', 'categories', 'orders', 'items'];
      for (const table of commonTables) {
        try {
          const { data, error } = await supabase.from(table).select('*').limit(1);
          if (!error) {
            tables.push({ table_name: table });
          }
        } catch (e) {
          // Table doesn't exist, ignore
        }
      }
      
      return res.json(tables);
    }
    
    return res.json(data);
  } catch (error) {
    console.error('Error listing tables:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Create a new table
 */
async function createTable(req, res, supabase) {
  try {
    const { schema = 'public', name, columns, enable_rls = false } = req.body;
    
    if (!name || !columns || !columns.length) {
      return res.status(400).json({ error: 'Table name and columns are required' });
    }
    
    const tableIdentifier = `"${schema}"."${name}"`;
    
    let sql = `CREATE TABLE ${tableIdentifier} (\n`;
    
    // Add columns
    const columnDefinitions = columns.map(column => {
      let def = `  "${column.name}" ${column.type}`;
      
      // Not null
      if (column.nullable === false) {
        def += ' NOT NULL';
      }
      
      // Default value
      if (column.default_value) {
        def += ` DEFAULT ${column.default_value}`;
      }
      
      // Primary key
      if (column.primary_key) {
        def += ' PRIMARY KEY';
      }
      
      // Unique
      if (column.unique) {
        def += ' UNIQUE';
      }
      
      // References (foreign key)
      if (column.references) {
        def += ` REFERENCES "${column.references.table}" ("${column.references.column}")`;
      }
      
      return def;
    });
    
    sql += columnDefinitions.join(',\n');
    sql += '\n)';
    
    // Enable RLS if requested
    if (enable_rls) {
      sql += `;\nALTER TABLE ${tableIdentifier} ENABLE ROW LEVEL SECURITY`;
    }
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Table '${name}' created successfully in schema '${schema}'`
    });
  } catch (error) {
    console.error('Error creating table:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Alter an existing table
 */
async function alterTable(req, res, supabase) {
  try {
    const { 
      schema = 'public', 
      name, 
      new_name, 
      add_columns = [], 
      drop_columns = [],
      enable_rls
    } = req.body;
    
    if (!name) {
      return res.status(400).json({ error: 'Table name is required' });
    }
    
    const tableIdentifier = `"${schema}"."${name}"`;
    let sqls = [];
    
    // Rename table
    if (new_name) {
      sqls.push(`ALTER TABLE ${tableIdentifier} RENAME TO "${new_name}"`);
    }
    
    // Add columns
    if (add_columns.length > 0) {
      for (const column of add_columns) {
        let sql = `ALTER TABLE ${tableIdentifier} ADD COLUMN "${column.name}" ${column.type}`;
        
        if (column.nullable === false) {
          sql += ' NOT NULL';
        }
        
        if (column.default_value) {
          sql += ` DEFAULT ${column.default_value}`;
        }
        
        sqls.push(sql);
      }
    }
    
    // Drop columns
    if (drop_columns.length > 0) {
      for (const columnName of drop_columns) {
        sqls.push(`ALTER TABLE ${tableIdentifier} DROP COLUMN "${columnName}"`);
      }
    }
    
    // Enable/disable RLS
    if (enable_rls !== undefined) {
      if (enable_rls) {
        sqls.push(`ALTER TABLE ${tableIdentifier} ENABLE ROW LEVEL SECURITY`);
      } else {
        sqls.push(`ALTER TABLE ${tableIdentifier} DISABLE ROW LEVEL SECURITY`);
      }
    }
    
    // Execute all queries
    const results = [];
    
    for (const sql of sqls) {
      const { data, error } = await supabase.rpc('execute_sql', { query: sql });
      
      if (error) {
        results.push({ 
          success: false, 
          query: sql, 
          error: error.message 
        });
      } else {
        results.push({ 
          success: true, 
          query: sql 
        });
      }
    }
    
    // Check if all operations were successful
    const allSuccess = results.every(r => r.success);
    
    return res.json({ 
      success: allSuccess, 
      operations: results,
      message: allSuccess 
        ? `Table '${name}' altered successfully` 
        : `Some operations failed while altering table '${name}'`
    });
  } catch (error) {
    console.error('Error altering table:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Drop a table
 */
async function dropTable(req, res, supabase) {
  try {
    const { schema = 'public', name, cascade = false } = req.body;
    
    if (!name) {
      return res.status(400).json({ error: 'Table name is required' });
    }
    
    const tableIdentifier = `"${schema}"."${name}"`;
    let sql = `DROP TABLE IF EXISTS ${tableIdentifier}`;
    
    if (cascade) {
      sql += ` CASCADE`;
    }
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Table '${name}' dropped successfully from schema '${schema}'`
    });
  } catch (error) {
    console.error('Error dropping table:', error);
    return res.status(500).json({ error: error.message });
  }
}

module.exports = {
  queryTable,
  generateTypes,
  listTables,
  createTable,
  alterTable,
  dropTable
};
