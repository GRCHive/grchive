package pg_api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strings"
	"time"
)

type PgDriver struct {
	connInfo   *core.DatabaseConnection
	db         *sqlx.DB
	connection *sql.Conn
	state      *core.FullDbState
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

	// Force a timeout
	build.WriteString(fmt.Sprintf("%s=%d", "connect_timeout", 10))
	return build.String()
}

func (pg *PgDriver) Connect(conn *core.DatabaseConnection, ro bool) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	pg.connInfo = conn
	pg.db, err = sqlx.ConnectContext(ctx, "postgres", CreateDatabaseConnectionString(conn))
	if err != nil {
		return err
	}

	pg.connection, err = pg.db.Conn(context.Background())
	if err != nil {
		return err
	}

	if ro {
		_, err := pg.connection.ExecContext(context.Background(), `
			SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY
		`)

		if err != nil {
			return err
		}
	}

	return nil
}
