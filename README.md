# Servidor MCP para Supabase Self-Hosted

Este é um servidor MCP (Model Context Protocol) personalizado para conectar ferramentas de IA (como Claude, Cursor, etc.) com instalações self-hosted do Supabase.

## Características

- Consulta de tabelas no banco de dados PostgreSQL do Supabase
- Listagem de tabelas disponíveis
- Geração de tipos TypeScript para facilitar o desenvolvimento
- Execução de consultas SQL com restrições de segurança
- Compatível com o protocolo MCP para integração com ferramentas de IA
- Gerenciamento de RLS (Row Level Security)
- Gerenciamento de Edge Functions
- Gerenciamento de esquemas de banco de dados
- Gerenciamento de tabelas
- Gerenciamento de buckets de armazenamento e suas políticas

## Requisitos

- Go 1.21 ou superior
- Docker (opcional, para execução em contêiner)
- Acesso a uma instalação Supabase self-hosted
- Service Role Key do seu projeto Supabase
- Docker (para execução containerizada)

## Instalação

### Instalação via Docker (Recomendado)

1. Use a imagem Docker:
   ```bash
   docker run -p 3000:3000 -e SUPABASE_URL=http://seu-servidor-supabase:8000 -e SUPABASE_KEY=seu-service-role-key dirgocs/supabase-self-hosted-mcp
   ```

### Instalação Manual

1. Clone este repositório:
   ```bash
   git clone https://github.com/dirgocs/supabase-self-hosted-mcp.git
   cd supabase-self-hosted-mcp
   ```

2. Instale as dependências:
   ```bash
   npm install
   ```

3. Configure as variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
   
   Edite o arquivo `.env` e configure as seguintes variáveis:
   - `SUPABASE_URL`: URL da sua instalação Supabase self-hosted
   - `SUPABASE_KEY`: Service Role Key do seu projeto Supabase
   - `PORT`: Porta em que o servidor MCP será executado (padrão: 3000)

4. Execute o servidor:
   ```bash
   npm start
   ```

## Uso com Claude ou outras ferramentas de IA

1. Inicie o servidor MCP usando o comando acima
2. Configure sua ferramenta de IA para usar o servidor MCP:

### Configuração para Claude

```json
{
  "mcpServers": {
    "supabase-self-hosted": {
      "command": "node",
      "args": ["index.js"],
      "cwd": "/caminho/para/supabase-self-hosted-mcp",
      "env": {
        "SUPABASE_URL": "http://seu-servidor-supabase:8000",
        "SUPABASE_KEY": "seu-service-role-key"
      }
    }
  }
}
```

### Configuração para Cursor ou Windsurf

```json
{
  "name": "Supabase Self-Hosted",
  "command": "node",
  "args": ["index.js"],
  "cwd": "/caminho/para/supabase-self-hosted-mcp",
  "env": {
    "SUPABASE_URL": "http://seu-servidor-supabase:8000", 
    "SUPABASE_KEY": "seu-service-role-key"
  }
}
```

## Funções disponíveis

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
