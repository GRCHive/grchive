package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/grchive/grchive/core"
	"time"
)

var dbConn *sqlx.DB

func Init() {
	dbConn = sqlx.MustConnect("postgres", core.EnvConfig.DatabaseConnString)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(5 * time.Minute)
}

func CreateTx() *sqlx.Tx {
	return dbConn.MustBegin()
}
