package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

var dbConn *sqlx.DB

func Init() {
	envConfig := core.LoadEnvConfig()

	dbConn = sqlx.MustConnect(envConfig.DatabaseDriver, envConfig.DatabaseConnString)
}
