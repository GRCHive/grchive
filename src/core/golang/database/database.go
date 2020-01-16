package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

var dbConn *sqlx.DB

func Init() {
	dbConn = sqlx.MustConnect("postgres", core.EnvConfig.DatabaseConnString)
}

func CreateTx() *sqlx.Tx {
	return dbConn.MustBegin()
}
