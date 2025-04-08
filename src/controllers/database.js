// Controllers for database-related operations
const { isReadOnlyQuery } = require('../utils');

/**
 * Execute a SQL query (read-only for security)
 */
async function executeQuery(req, res, supabase) {
  try {
    const { query } = req.body;
    
    if (!query) {
      return res.status(400).json({ error: 'Query is required' });
    }
    
    // Check if the query is read-only
    if (!isReadOnlyQuery(query)) {
      return res.status(403).json({ 
        error: 'Only read-only queries are allowed through this endpoint for security reasons'
      });
    }
    
    // Execute query through RPC
    const { data, error } = await supabase.rpc('execute_sql', { query });
    
    if (error) {
      return res.status(400).json({ 
        error: error.message,
        message: "Unable to execute query. You may need to create a custom function 'execute_sql' in your Supabase instance."
      });
    }
    
    return res.json(data);
  } catch (error) {
    console.error('Error executing query:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Get database schema information
 */
async function getDatabaseSchema(req, res, supabase) {
  try {
    const { schema } = req.body;
    
    let query = `
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
    `;
    
    if (schema) {
      query += ` AND n.nspname = '${schema}'`;
    } else {
      query += ` AND n.nspname NOT IN ('pg_catalog', 'information_schema')`;
    }
    
    query += ` ORDER BY n.nspname, c.relname, a.attnum`;
    
    const { data, error } = await supabase.rpc('execute_sql', { query });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    // Process data into a more user-friendly structure
    const schemaData = {};
    
    (data || []).forEach(item => {
      const schemaName = item.schema_name;
      const tableName = item.table_name;
      
      if (!schemaData[schemaName]) {
        schemaData[schemaName] = {};
      }
      
      if (!schemaData[schemaName][tableName]) {
        schemaData[schemaName][tableName] = {
          columns: [],
          primary_keys: [],
          foreign_keys: []
        };
      }
      
      // Add column
      const column = {
        name: item.column_name,
        type: item.data_type,
        not_null: item.not_null,
        default_value: item.default_value
      };
      
      schemaData[schemaName][tableName].columns.push(column);
      
      // Add primary key
      if (item.is_primary_key) {
        schemaData[schemaName][tableName].primary_keys.push(item.column_name);
      }
      
      // Add foreign key
      if (item.is_foreign_key) {
        schemaData[schemaName][tableName].foreign_keys.push({
          column: item.column_name,
          references: {
            schema: item.reference_schema,
            table: item.reference_table,
            column: item.reference_column
          }
        });
      }
    });
    
    return res.json(schemaData);
  } catch (error) {
    console.error('Error getting database schema:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Create a new schema
 */
async function createSchema(req, res, supabase) {
  try {
    const { name } = req.body;
    
    if (!name) {
      return res.status(400).json({ error: 'Schema name is required' });
    }
    
    const sql = `CREATE SCHEMA IF NOT EXISTS "${name}"`;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Schema '${name}' created successfully`
    });
  } catch (error) {
    console.error('Error creating schema:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Delete a schema
 */
async function deleteSchema(req, res, supabase) {
  try {
    const { name, cascade = false } = req.body;
    
    if (!name) {
      return res.status(400).json({ error: 'Schema name is required' });
    }
    
    let sql = `DROP SCHEMA IF EXISTS "${name}"`;
    
    if (cascade) {
      sql += ` CASCADE`;
    }
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Schema '${name}' deleted successfully`
    });
  } catch (error) {
    console.error('Error deleting schema:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Get RLS policies
 */
async function getRlsPolicies(req, res, supabase) {
  try {
    const { schema = 'public', table } = req.body;
    
    let query = `
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
    `;
    
    const params = [schema];
    
    if (table) {
      query += ` AND c.relname = $2`;
      params.push(table);
    }
    
    query += ` ORDER BY n.nspname, c.relname, p.polname`;
    
    const { data, error } = await supabase.rpc('execute_sql', { 
      query,
      params
    });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json(data || []);
  } catch (error) {
    console.error('Error getting RLS policies:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Create a new RLS policy
 */
async function createRlsPolicy(req, res, supabase) {
  try {
    const { 
      schema = 'public', 
      table, 
      name, 
      operation, 
      definition, 
      check = null, 
      role = 'public' 
    } = req.body;
    
    if (!table || !name || !operation || !definition) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    const tableIdentifier = `"${schema}"."${table}"`;
    
    let sql = `CREATE POLICY "${name}" ON ${tableIdentifier} 
               FOR ${operation} 
               TO "${role}" 
               USING (${definition})`;
               
    // Add WITH CHECK clause for operations that need it
    if (check && (operation === 'INSERT' || operation === 'UPDATE' || operation === 'ALL')) {
      sql += ` WITH CHECK (${check})`;
    }
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `RLS policy '${name}' created on ${tableIdentifier}`
    });
  } catch (error) {
    console.error('Error creating RLS policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Update an existing RLS policy
 */
async function updateRlsPolicy(req, res, supabase) {
  try {
    const { 
      schema = 'public', 
      table, 
      name, 
      operation = null, 
      definition,
      check = null
    } = req.body;
    
    if (!table || !name || !definition) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    const tableIdentifier = `"${schema}"."${table}"`;
    
    // First drop the existing policy
    const dropSql = `DROP POLICY IF EXISTS "${name}" ON ${tableIdentifier}`;
    await supabase.rpc('execute_sql', { query: dropSql });
    
    // Then recreate it with new parameters
    const createSql = `CREATE POLICY "${name}" ON ${tableIdentifier} 
                      FOR ${operation || 'ALL'} 
                      USING (${definition})`;
                       
    // Add WITH CHECK clause for operations that need it
    const createWithCheck = operation && check && 
      (operation === 'INSERT' || operation === 'UPDATE' || operation === 'ALL') 
        ? ` WITH CHECK (${check})` 
        : '';
        
    const finalCreateSql = createSql + createWithCheck;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: finalCreateSql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `RLS policy '${name}' updated on ${tableIdentifier}`
    });
  } catch (error) {
    console.error('Error updating RLS policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Delete an RLS policy
 */
async function deleteRlsPolicy(req, res, supabase) {
  try {
    const { schema = 'public', table, name } = req.body;
    
    if (!table || !name) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    const tableIdentifier = `"${schema}"."${table}"`;
    const sql = `DROP POLICY IF EXISTS "${name}" ON ${tableIdentifier}`;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `RLS policy '${name}' deleted from ${tableIdentifier}`
    });
  } catch (error) {
    console.error('Error deleting RLS policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

module.exports = {
  executeQuery,
  getDatabaseSchema,
  createSchema,
  deleteSchema,
  getRlsPolicies,
  createRlsPolicy,
  updateRlsPolicy,
  deleteRlsPolicy
};
