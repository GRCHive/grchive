package db_api

import (
	"gitlab.com/grchive/grchive/core"
)

type DbDriver interface {
	Connect(*core.DatabaseConnection) error
	ConnectionReadOnly() bool

	GetSchemas() ([]*core.DbSchema, error)
	GetTables(*core.DbSchema) ([]*core.DbTable, error)
	GetColumns(*core.DbTable) ([]*core.DbColumn, error)
}
