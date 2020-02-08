package pg_api

import (
	"gitlab.com/grchive/grchive/core"
)

type PgDriver struct {
}

func (pg *PgDriver) Connect(conn *core.DatabaseConnection) error {
	return nil
}
