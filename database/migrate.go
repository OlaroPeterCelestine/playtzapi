package database

import (
	"database/sql"
	"embed"
)

//go:embed schema.sql
var schemaFS embed.FS

// Migrate runs the database migrations
func Migrate() error {
	if DB == nil {
		return sql.ErrConnDone
	}

	// Read the schema file
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	// Execute the schema
	_, err = DB.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
