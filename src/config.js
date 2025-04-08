// Configuration for Supabase Self-Hosted MCP Server
require('dotenv').config();

// Supabase configuration
const supabaseConfig = {
  url: process.env.SUPABASE_URL || 'http://your-self-hosted-supabase-url',
  key: process.env.SUPABASE_KEY || 'your-supabase-key'
};

// Server configuration
const serverConfig = {
  port: process.env.PORT || 3000,
  env: process.env.NODE_ENV || 'development'
};

module.exports = {
  supabase: supabaseConfig,
  server: serverConfig
};
