package pg_api

import (
	"gitlab.com/grchive/grchive/core"
	"strings"
)

func (pg *PgDriver) GetSchemas() ([]*core.DbSchema, error) {
	rows, err := pg.connection.Queryx(`
		SELECT schema_name
		FROM information_schema.schemata
		WHERE catalog_name = $1
	`, pg.connInfo.DbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schemas := make([]*core.DbSchema, 0)

	for rows.Next() {
		s := core.DbSchema{}
		err = rows.Scan(&s.SchemaName)
		if err != nil {
			return nil, err
		}

		// Ignore default schemas provided by the Postgres database.
		if s.SchemaName == "information_schema" ||
			s.SchemaName == "pg_catalog" ||
			strings.HasPrefix(s.SchemaName, "pg_toast") ||
			strings.HasPrefix(s.SchemaName, "pg_temp") {
			continue
		}

		schemas = append(schemas, &s)
	}

	return schemas, nil
}

func (pg *PgDriver) GetTables(schema *core.DbSchema) ([]*core.DbTable, error) {
	rows, err := pg.connection.Queryx(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_catalog = $1
			AND table_schema = $2
	`, pg.connInfo.DbName, schema.SchemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]*core.DbTable, 0)

	for rows.Next() {
		s := core.DbTable{}
		err = rows.Scan(&s.TableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &s)
	}

	return tables, nil
}

func (pg *PgDriver) GetColumns(schema *core.DbSchema, table *core.DbTable) ([]*core.DbColumn, error) {
	rows, err := pg.connection.Queryx(`
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_catalog = $1
			AND table_schema = $2
			AND table_name = $3
	`, pg.connInfo.DbName, schema.SchemaName, table.TableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]*core.DbColumn, 0)

	for rows.Next() {
		s := core.DbColumn{}
		err = rows.Scan(&s.ColumnName, &s.ColumnType)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &s)
	}

	return columns, nil
}
