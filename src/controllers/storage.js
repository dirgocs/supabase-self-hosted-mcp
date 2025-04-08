// Controllers for storage-related operations

/**
 * Get storage buckets
 */
async function getBuckets(req, res, supabase) {
  try {
    const { id } = req.body;
    
    let query = `
      SELECT * FROM storage.buckets
    `;
    
    if (id) {
      query += ` WHERE id = '${id}'`;
    }
    
    const { data, error } = await supabase.rpc('execute_sql', { query });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json(data || []);
  } catch (error) {
    console.error('Error getting buckets:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Create a new storage bucket
 */
async function createBucket(req, res, supabase) {
  try {
    const { 
      id, 
      name = null, 
      public = false,
      file_size_limit = null,
      allowed_mime_types = null
    } = req.body;
    
    if (!id) {
      return res.status(400).json({ error: 'Bucket ID is required' });
    }
    
    // Try using Supabase API first
    try {
      const { data, error } = await supabase
        .storage
        .createBucket(id, {
          public,
          fileSizeLimit: file_size_limit,
          allowedMimeTypes: allowed_mime_types
        });
      
      if (!error) {
        return res.json({ 
          success: true, 
          message: `Bucket '${id}' created successfully`,
          method: 'api',
          data
        });
      }
    } catch (apiError) {
      console.log('Storage API not available, falling back to SQL');
    }
    
    // Fallback to SQL
    let sql = `
      INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types, created_at, updated_at)
      VALUES (
        '${id}',
        '${name || id}',
        ${public},
        ${file_size_limit ? file_size_limit : 'NULL'},
        ${allowed_mime_types ? `'${JSON.stringify(allowed_mime_types)}'::jsonb` : 'NULL'},
        NOW(),
        NOW()
      )
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Bucket '${id}' created successfully`,
      method: 'sql'
    });
  } catch (error) {
    console.error('Error creating bucket:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Update a storage bucket
 */
async function updateBucket(req, res, supabase) {
  try {
    const { 
      id, 
      public,
      file_size_limit,
      allowed_mime_types
    } = req.body;
    
    if (!id) {
      return res.status(400).json({ error: 'Bucket ID is required' });
    }
    
    // Try using Supabase API first
    try {
      const { data, error } = await supabase
        .storage
        .updateBucket(id, {
          public,
          fileSizeLimit: file_size_limit,
          allowedMimeTypes: allowed_mime_types
        });
      
      if (!error) {
        return res.json({ 
          success: true, 
          message: `Bucket '${id}' updated successfully`,
          method: 'api',
          data
        });
      }
    } catch (apiError) {
      console.log('Storage API not available, falling back to SQL');
    }
    
    // Fallback to SQL
    let sql = `UPDATE storage.buckets SET updated_at = NOW()`;
    
    if (public !== undefined) {
      sql += `, public = ${public}`;
    }
    
    if (file_size_limit !== undefined) {
      sql += `, file_size_limit = ${file_size_limit || 'NULL'}`;
    }
    
    if (allowed_mime_types !== undefined) {
      sql += `, allowed_mime_types = ${allowed_mime_types ? `'${JSON.stringify(allowed_mime_types)}'::jsonb` : 'NULL'}`;
    }
    
    sql += ` WHERE id = '${id}'`;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Bucket '${id}' updated successfully`,
      method: 'sql'
    });
  } catch (error) {
    console.error('Error updating bucket:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Delete a storage bucket
 */
async function deleteBucket(req, res, supabase) {
  try {
    const { id } = req.body;
    
    if (!id) {
      return res.status(400).json({ error: 'Bucket ID is required' });
    }
    
    // Try using Supabase API first
    try {
      const { data, error } = await supabase
        .storage
        .deleteBucket(id);
      
      if (!error) {
        return res.json({ 
          success: true, 
          message: `Bucket '${id}' deleted successfully`,
          method: 'api'
        });
      }
    } catch (apiError) {
      console.log('Storage API not available, falling back to SQL');
    }
    
    // Fallback to SQL
    const sql = `DELETE FROM storage.buckets WHERE id = '${id}'`;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Bucket '${id}' deleted successfully`,
      method: 'sql'
    });
  } catch (error) {
    console.error('Error deleting bucket:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Get bucket policies
 */
async function getBucketPolicies(req, res, supabase) {
  try {
    const { bucket_id } = req.body;
    
    if (!bucket_id) {
      return res.status(400).json({ error: 'Bucket ID is required' });
    }
    
    const query = `
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
      WHERE name = '${bucket_id}'
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json(data || []);
  } catch (error) {
    console.error('Error getting bucket policies:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Create a bucket policy
 */
async function createBucketPolicy(req, res, supabase) {
  try {
    const { 
      bucket_id, 
      name, 
      operation, 
      definition, 
      role = 'authenticated'
    } = req.body;
    
    if (!bucket_id || !name || !operation || !definition) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    // Map operation string to internal code
    const opMap = {
      'SELECT': 10,
      'INSERT': 20,
      'UPDATE': 40,
      'DELETE': 80
    };
    
    const opCode = opMap[operation];
    
    if (!opCode) {
      return res.status(400).json({ error: 'Invalid operation. Must be SELECT, INSERT, UPDATE, or DELETE' });
    }
    
    const sql = `
      INSERT INTO storage.policies (name, bucket_id, operation, definition, role, created_at)
      VALUES (
        '${name}',
        '${bucket_id}',
        ${opCode},
        '${definition.replace(/'/g, "''")}',
        '${role}',
        NOW()
      )
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Policy '${name}' for bucket '${bucket_id}' created successfully`
    });
  } catch (error) {
    console.error('Error creating bucket policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Update a bucket policy
 */
async function updateBucketPolicy(req, res, supabase) {
  try {
    const { bucket_id, name, definition } = req.body;
    
    if (!bucket_id || !name || !definition) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    const sql = `
      UPDATE storage.policies
      SET 
        definition = '${definition.replace(/'/g, "''")}',
        updated_at = NOW()
      WHERE 
        bucket_id = '${bucket_id}' AND 
        name = '${name}'
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Policy '${name}' for bucket '${bucket_id}' updated successfully`
    });
  } catch (error) {
    console.error('Error updating bucket policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

/**
 * Delete a bucket policy
 */
async function deleteBucketPolicy(req, res, supabase) {
  try {
    const { bucket_id, name } = req.body;
    
    if (!bucket_id || !name) {
      return res.status(400).json({ error: 'Missing required parameters' });
    }
    
    const sql = `
      DELETE FROM storage.policies
      WHERE 
        bucket_id = '${bucket_id}' AND 
        name = '${name}'
    `;
    
    const { data, error } = await supabase.rpc('execute_sql', { query: sql });
    
    if (error) {
      return res.status(400).json({ error: error.message });
    }
    
    return res.json({ 
      success: true, 
      message: `Policy '${name}' for bucket '${bucket_id}' deleted successfully`
    });
  } catch (error) {
    console.error('Error deleting bucket policy:', error);
    return res.status(500).json({ error: error.message });
  }
}

module.exports = {
  getBuckets,
  createBucket,
  updateBucket,
  deleteBucket,
  getBucketPolicies,
  createBucketPolicy,
  updateBucketPolicy,
  deleteBucketPolicy
};
