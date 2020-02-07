package db_api

import (
	"gitlab.com/grchive/grchive/core"
)

func CreateDriver(dbType *core.DatabaseType, dbConn *core.DatabaseConnection) (DbDriver, error) {
	var driver DbDriver
	driver = nil
	return driver, nil
}
