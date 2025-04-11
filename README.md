# Supabase Self-Hosted MCP Server

This is a custom MCP (Model Context Protocol) server implemented in Go for connecting AI tools (like Claude, Cursor, etc.) with self-hosted Supabase installations. It provides a comprehensive API for managing your Supabase resources.

## Features

- Query tables in the Supabase PostgreSQL database
- List available tables and schemas
- Generate TypeScript types to facilitate development
- Execute SQL queries with security restrictions (including migrations)
- Compatible with the MCP protocol for integration with AI tools
- Row Level Security (RLS) management
- Edge Functions management
- Database schema management
- Table management
- Storage bucket management and policies
- RESTful API for programmatic access

## Requirements

- Go 1.21 or higher
- Docker (optional, for containerized execution)
- Access to a self-hosted Supabase installation
- Service Role Key from your Supabase project
- Coolify v4.0.0-beta.406 or higher (for Coolify deployment)

## Installation

### Docker Installation (Recommended)

1. Use the Docker image:
   ```bash
   docker run -p 3000:3000 -e SUPABASE_URL=http://your-supabase-server:8000 -e SUPABASE_KEY=your-service-role-key dirgocs/supabase-self-hosted-mcp:1.0.0
   ```

### Manual Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/dirgocs/supabase-self-hosted-mcp.git
   cd supabase-self-hosted-mcp
   ```

2. Build the application:
   ```bash
   go build -o supabase-mcp
   ```

3. Configure environment variables:
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file and configure the following variables:
   - `SUPABASE_URL`: URL of your self-hosted Supabase installation
   - `SUPABASE_KEY`: Service Role Key from your Supabase project
   - `SUPABASE_ANON_KEY`: Anonymous Key for public operations
   - `SUPABASE_JWT_SECRET`: JWT Secret used for token verification
   - `PG_CONNECTION_STRING`: Direct PostgreSQL connection string (optional)
   - `PORT`: Port on which the MCP server will run (default: 3000)
   - `GIN_MODE`: Gin framework mode (debug or release)

4. Run the server:
   ```bash
   ./supabase-mcp
   ```

### Coolify Installation

This MCP server is designed to work with Coolify v4.0.0-beta.406 or higher. You can deploy it directly from the GitHub repository using Coolify's Docker deployment feature.

1. In Coolify, create a new service
2. Select GitHub repository as the source
3. Configure the environment variables:
   - `SUPABASE_URL`: URL of your self-hosted Supabase installation
   - `SUPABASE_KEY`: Service Role Key from your Supabase project
   - `SUPABASE_ANON_KEY`: Anonymous Key for public operations
   - `SUPABASE_JWT_SECRET`: JWT Secret for token verification
   - `PG_CONNECTION_STRING`: Direct PostgreSQL connection string (optional)
4. Deploy the service

The Docker configuration is already optimized for Coolify deployment.

## Using with Claude or other AI tools

1. Start the MCP server using the command above
2. Configure your AI tool to use the MCP server:

### Configuration for Claude

```json
{
  "mcpServers": {
    "supabase-self-hosted": {
      "command": "./supabase-mcp",
      "cwd": "/path/to/supabase-self-hosted-mcp",
      "env": {
        "SUPABASE_URL": "http://your-supabase-server:8000",
        "SUPABASE_KEY": "your-service-role-key"
      }
    }
  }
}
```

### Configuration for Cursor or Windsurf

```json
{
  "name": "Supabase Self-Hosted",
  "command": "./supabase-mcp",
  "cwd": "/path/to/supabase-self-hosted-mcp",
  "env": {
    "SUPABASE_URL": "http://your-supabase-server:8000", 
    "SUPABASE_KEY": "your-service-role-key"
  }
}
```

### Using with Docker

If you're using the Docker container, you can configure your AI tool to connect to the exposed port (default: 3000) of the container.

## Available Endpoints

The MCP server provides the following API endpoints:

### Database Management
- `GET /api/database/schema`: Get database schema information
- `POST /api/database/query`: Execute SQL queries
- `GET /api/database/rls`: Get RLS policies

### Table Management
- `GET /api/tables`: List available tables
- `GET /api/tables/:table`: Get table details
- `POST /api/tables/:table/query`: Query table data
- `GET /api/tables/types`: Generate TypeScript types

### Storage Management
- `GET /api/storage/buckets`: List storage buckets
- `POST /api/storage/buckets`: Create a new bucket
- `PUT /api/storage/buckets/:id`: Update bucket settings
- `DELETE /api/storage/buckets/:id`: Delete a bucket

### Edge Functions Management
- `GET /api/edge-functions`: List edge functions
- `POST /api/edge-functions`: Create a new edge function
- `PUT /api/edge-functions/:name`: Update an edge function
- `DELETE /api/edge-functions/:name`: Delete an edge function

## Running SQL Migrations

You can use the `/api/database/query` endpoint to run SQL migrations. Simply send a POST request with your SQL migration script:

```bash
curl -X POST http://localhost:3000/api/database/query \
  -H "Content-Type: application/json" \
  -d '{"query": "CREATE TABLE IF NOT EXISTS my_table (id SERIAL PRIMARY KEY, name TEXT);"}'
```

For complex migrations, you can execute multiple statements in a transaction:

```bash
curl -X POST http://localhost:3000/api/database/query \
  -H "Content-Type: application/json" \
  -d '{"query": "BEGIN; CREATE TABLE IF NOT EXISTS my_table (id SERIAL PRIMARY KEY, name TEXT); CREATE INDEX idx_my_table_name ON my_table(name); COMMIT;"}'
```

## Docker Network Configuration

When running with Docker, you can use a shared network to connect to your Supabase services:

```yaml
services:
  mcp-server:
    image: dirgocs/supabase-self-hosted-mcp:1.0.0
    environment:
      - SUPABASE_URL=http://kong:8000
      - SUPABASE_KEY=your-service-role-key
    networks:
      - supabase-network

networks:
  supabase-network:
    external: true
```

## Security Considerations

This MCP server uses the Supabase Service Role Key, which has full access to your database. Make sure to:

1. Keep your `.env` file secure and never commit it to version control
2. Restrict access to the MCP server to trusted users only
3. Consider implementing additional authentication for the MCP server
4. Run the server in a secure environment

### Tabelas e Consultas
- `query_table`: Consultar uma tabela específica com suporte a filtros
- `generate_types`: Gerar tipos TypeScript para seu esquema de banco de dados
- `list_tables`: Listar todas as tabelas em um esquema específico
- `execute_query`: Executar uma consulta SQL (apenas operações de leitura)

### Row Level Security (RLS)
- `get_rls_policies`: Obter políticas RLS para uma tabela ou todas as tabelas
- `create_rls_policy`: Criar uma nova política RLS
- `update_rls_policy`: Atualizar uma política RLS existente
- `delete_rls_policy`: Excluir uma política RLS

### Edge Functions
- `get_edge_functions`: Obter todas as edge functions ou uma específica
- `create_edge_function`: Criar uma nova edge function
- `update_edge_function`: Atualizar uma edge function existente
- `delete_edge_function`: Excluir uma edge function
- `deploy_edge_function`: Implantar uma edge function

### Esquema de Banco de Dados
- `get_database_schema`: Obter o esquema do banco de dados
- `create_schema`: Criar um novo esquema
- `delete_schema`: Excluir um esquema

### Tabelas
- `create_table`: Criar uma nova tabela
- `alter_table`: Alterar uma tabela (adicionar/remover colunas, renomear)
- `drop_table`: Excluir uma tabela

### Buckets de Armazenamento
- `get_buckets`: Obter todos os buckets de armazenamento ou um específico
- `create_bucket`: Criar um novo bucket de armazenamento
- `update_bucket`: Atualizar um bucket de armazenamento
- `delete_bucket`: Excluir um bucket de armazenamento

### Políticas de Bucket
- `get_bucket_policies`: Obter políticas para um bucket de armazenamento
- `create_bucket_policy`: Criar uma nova política para um bucket
- `update_bucket_policy`: Atualizar uma política de bucket
- `delete_bucket_policy`: Excluir uma política de bucket

## Segurança

Este servidor implementa as seguintes medidas de segurança:

- Restrição a operações apenas de leitura no endpoint `execute_query`
- Uso da API oficial do Supabase para consultas
- Validação de parâmetros de entrada

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## Licença

MIT
