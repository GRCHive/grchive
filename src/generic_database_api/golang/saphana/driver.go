package saphana_api

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/SAP/go-hdb/driver"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strings"
	"time"
)

type SapHanaDriver struct {
	connInfo   *core.DatabaseConnection
	db         *sqlx.DB
	connection *sql.Conn
	state      *core.FullDbState
}

func CreateDatabaseConnectionString(conn *core.DatabaseConnection) string {
	build := strings.Builder{}
	_, err := build.WriteString(fmt.Sprintf("sap://%s:%s@%s:%d?",
		conn.Username,
		conn.Password,
		conn.Host,
		conn.Port,
	))
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

func (hdb *SapHanaDriver) Connect(conn *core.DatabaseConnection, ro bool) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	hdb.connInfo = conn
	hdb.db, err = sqlx.ConnectContext(ctx, "hdb", CreateDatabaseConnectionString(conn))
	if err != nil {
		return err
	}

	hdb.connection, err = hdb.db.Conn(context.Background())
	if err != nil {
		return err
	}

	if ro {
		_, err = hdb.connection.ExecContext(context.Background(), `
			SET TRANSACTION READ ONLY
		`)
		if err != nil {
			return err
		}
	}

	return nil
}
