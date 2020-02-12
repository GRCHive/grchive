package pg_api

import (
	"gitlab.com/grchive/grchive/core"
)

func (pg *PgDriver) GetSchemas() ([]*core.DbSchema, error) {
	return nil, nil
}

func (pg *PgDriver) GetTables(schema *core.DbSchema) ([]*core.DbTable, error) {
	return nil, nil
}

func (pg *PgDriver) GetColumns(table *core.DbTable) ([]*core.DbColumn, error) {
	return nil, nil
}
