// Main MCP Server for Supabase Self-Hosted
const express = require('express');
const { createClient } = require('@supabase/supabase-js');
const cors = require('cors');
const bodyParser = require('body-parser');
const morgan = require('morgan');
require('dotenv').config();

// Import configuration and specifications
const config = require('./src/config');
const mcpSpec = require('./src/mcp-spec');

// Import controllers
const tableControllers = require('./src/controllers/tables');
const storageControllers = require('./src/controllers/storage');
const databaseControllers = require('./src/controllers/database');
const edgeFunctionControllers = require('./src/controllers/edge-functions');

// Initialize Express app
const app = express();
const PORT = config.server.port;

// Apply middleware
app.use(cors());
app.use(bodyParser.json());
app.use(morgan('dev'));

// Create Supabase client
const supabase = createClient(config.supabase.url, config.supabase.key);

// MCP Server info endpoint
app.get('/', (req, res) => {
  res.json({
    name: 'Supabase Self-Hosted MCP Server',
    version: '1.0.0',
    description: 'MCP Server para comunicação com Supabase Self-Hosted'
  });
});

// MCP specification endpoint
app.get('/v1/specification', (req, res) => {
  res.json(mcpSpec);
});

// Table endpoints
app.post('/v1/query_table', (req, res) => tableControllers.queryTable(req, res, supabase));
app.post('/v1/generate_types', (req, res) => tableControllers.generateTypes(req, res, supabase));
app.post('/v1/list_tables', (req, res) => tableControllers.listTables(req, res, supabase));
app.post('/v1/create_table', (req, res) => tableControllers.createTable(req, res, supabase));
app.post('/v1/alter_table', (req, res) => tableControllers.alterTable(req, res, supabase));
app.post('/v1/drop_table', (req, res) => tableControllers.dropTable(req, res, supabase));

// Database endpoints
app.post('/v1/execute_query', (req, res) => databaseControllers.executeQuery(req, res, supabase));
app.post('/v1/get_database_schema', (req, res) => databaseControllers.getDatabaseSchema(req, res, supabase));
app.post('/v1/create_schema', (req, res) => databaseControllers.createSchema(req, res, supabase));
app.post('/v1/delete_schema', (req, res) => databaseControllers.deleteSchema(req, res, supabase));

// RLS policy endpoints
app.post('/v1/get_rls_policies', (req, res) => databaseControllers.getRlsPolicies(req, res, supabase));
app.post('/v1/create_rls_policy', (req, res) => databaseControllers.createRlsPolicy(req, res, supabase));
app.post('/v1/update_rls_policy', (req, res) => databaseControllers.updateRlsPolicy(req, res, supabase));
app.post('/v1/delete_rls_policy', (req, res) => databaseControllers.deleteRlsPolicy(req, res, supabase));

// Edge function endpoints
app.post('/v1/get_edge_functions', (req, res) => edgeFunctionControllers.getEdgeFunctions(req, res, supabase));
app.post('/v1/create_edge_function', (req, res) => edgeFunctionControllers.createEdgeFunction(req, res, supabase));
app.post('/v1/update_edge_function', (req, res) => edgeFunctionControllers.updateEdgeFunction(req, res, supabase));
app.post('/v1/delete_edge_function', (req, res) => edgeFunctionControllers.deleteEdgeFunction(req, res, supabase));
app.post('/v1/deploy_edge_function', (req, res) => edgeFunctionControllers.deployEdgeFunction(req, res, supabase));

// Storage bucket endpoints
app.post('/v1/get_buckets', (req, res) => storageControllers.getBuckets(req, res, supabase));
app.post('/v1/create_bucket', (req, res) => storageControllers.createBucket(req, res, supabase));
app.post('/v1/update_bucket', (req, res) => storageControllers.updateBucket(req, res, supabase));
app.post('/v1/delete_bucket', (req, res) => storageControllers.deleteBucket(req, res, supabase));

// Bucket policy endpoints
app.post('/v1/get_bucket_policies', (req, res) => storageControllers.getBucketPolicies(req, res, supabase));
app.post('/v1/create_bucket_policy', (req, res) => storageControllers.createBucketPolicy(req, res, supabase));
app.post('/v1/update_bucket_policy', (req, res) => storageControllers.updateBucketPolicy(req, res, supabase));
app.post('/v1/delete_bucket_policy', (req, res) => storageControllers.deleteBucketPolicy(req, res, supabase));

// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Internal Server Error',
    message: err.message
  });
});

// Start the server
app.listen(PORT, () => {
  console.log(`Supabase Self-Hosted MCP Server running on port ${PORT}`);
  console.log(`Server URL: http://localhost:${PORT}`);
  console.log(`MCP Specification URL: http://localhost:${PORT}/v1/specification`);
  console.log(`Supabase URL: ${config.supabase.url}`);
});
