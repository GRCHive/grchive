package pg_api

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strings"
)

type PgDriver struct {
	connInfo   *core.DatabaseConnection
	connection *sqlx.DB

	currentRole *PgRole
	grants      *PgSchemaGrants
}

func CreateDatabaseConnectionString(conn *core.DatabaseConnection) string {
	build := strings.Builder{}
	_, err := build.WriteString(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?",
		conn.Username,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.DbName))
	if err != nil {
		core.Error("Failed to create connection string: " + err.Error())
	}

	for k, v := range conn.Parameters {
		_, err = build.WriteString(fmt.Sprintf("%s=%s&", k, v))
		if err != nil {
			core.Error("Failed to add parameters to connection string: " + err.Error())
		}
	}
	return build.String()
}

func (pg *PgDriver) Connect(conn *core.DatabaseConnection) error {
	var err error
	pg.connInfo = conn
	pg.connection, err = sqlx.Connect("postgres", CreateDatabaseConnectionString(conn))
	if err != nil {
		return err
	}

	pg.currentRole, err = retrieveRole(pg.connection, pg.connInfo.Username)
	if err != nil {
		return err
	}

	pg.grants, err = retrieveSchemaGrants(pg.connection, pg.connInfo.Username, pg.connInfo.DbName)
	if err != nil {
		return err
	}

	return nil
}
