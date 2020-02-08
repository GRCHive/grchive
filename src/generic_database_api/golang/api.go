package db_api

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/db_api/pg_api"
)

var POSTGRES_NAME string = "PostgreSQL"

func CreateDriver(dbType *core.DatabaseType, dbConn *core.DatabaseConnection) (DbDriver, error) {
	var driver DbDriver

	if dbType.Name == POSTGRES_NAME {
		driver = &pg_api.PgDriver{}
	} else {
		return nil, errors.New("Unsupported Database Type: " + dbType.Name)
	}

	if err := driver.Connect(dbConn); err != nil {
		return nil, err
	}

	return driver, nil
}
