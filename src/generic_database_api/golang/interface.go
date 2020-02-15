package db_api

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/db_api/utility"
)

type DbDriver interface {
	Connect(*core.DatabaseConnection) error
	ConnectionReadOnly() bool

	GetSchemas() ([]*core.DbSchema, error)
	GetTables(*core.DbSchema) ([]*core.DbTable, error)
	GetColumns(*core.DbSchema, *core.DbTable) ([]*core.DbColumn, error)
	GetFunctions(*core.DbSchema) ([]*core.DbFunction, error)

	RunQuery(query string) (*utility.SqlQueryResult, error)
}
