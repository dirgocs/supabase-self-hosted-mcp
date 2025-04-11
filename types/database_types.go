package types

// DatabaseColumn represents a database column
type DatabaseColumn struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	NotNull      bool        `json:"not_null"`
	DefaultValue interface{} `json:"default_value"`
}

// ForeignKeyReference represents a foreign key reference
type ForeignKeyReference struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	Column     string             `json:"column"`
	References ForeignKeyReference `json:"references"`
}

// TableSchema represents the schema of a database table
type TableSchema struct {
	Columns     []DatabaseColumn `json:"columns"`
	PrimaryKeys []string         `json:"primary_keys"`
	ForeignKeys []ForeignKey     `json:"foreign_keys"`
}

// SchemaData is a map of schema names to maps of table names to table schemas
type SchemaData map[string]map[string]TableSchema

// SchemaItem represents a row from the schema query result
type SchemaItem struct {
	SchemaName       string      `json:"schema_name"`
	TableName        string      `json:"table_name"`
	ColumnName       string      `json:"column_name"`
	DataType         string      `json:"data_type"`
	NotNull          bool        `json:"not_null"`
	DefaultValue     interface{} `json:"default_value"`
	IsPrimaryKey     bool        `json:"is_primary_key"`
	IsForeignKey     bool        `json:"is_foreign_key"`
	ReferencedTable  string      `json:"referenced_table"`
	ReferencedColumn string      `json:"referenced_column"`
}

// RLSPolicy represents a row-level security policy
type RLSPolicy struct {
	Schemaname string `json:"schemaname"`
	Tablename  string `json:"tablename"`
	Policyname string `json:"policyname"`
	Roles      string `json:"roles"`
	Cmd        string `json:"cmd"`
	Qual       string `json:"qual"`
	WithCheck  string `json:"with_check"`
}
