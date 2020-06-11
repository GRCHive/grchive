package db_api

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/db_api/utility"
)

type DbDriver interface {
	Connect(*core.DatabaseConnection, bool) error
	Close()

	ConstructState() error
	GetState() *core.FullDbState

	RunQuery(query string) (*utility.SqlQueryResult, error)
}
