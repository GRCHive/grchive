package pg_api

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

type PgDriver struct {
	connection *sqlx.DB
}

func CreateDatabaseConnectionString(conn *core.DatabaseConnection) string {
	baseStr := fmt.Sprintf("postgres://%s:%s@%s:%d", conn.Username, conn.Password, conn.Host, conn.Password)
	return baseStr
}

func (pg *PgDriver) Connect(conn *core.DatabaseConnection) error {
	var err error
	pg.connection, err = sqlx.Connect("postgres", CreateDatabaseConnectionString(conn))
	return err
}
