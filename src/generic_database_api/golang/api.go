package db_api

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/db_api/pg_api"
	"gitlab.com/grchive/grchive/db_api/saphana_api"
)

var POSTGRES_NAME string = "PostgreSQL"
var SAPHANA_NAME string = "SAP HANA"

func CreateDriver(dbType *core.DatabaseType, dbConn *core.DatabaseConnection, ro bool) (DbDriver, error) {
	var driver DbDriver

	if dbType.Name == POSTGRES_NAME {
		driver = &pg_api.PgDriver{}
	} else if dbType.Name == SAPHANA_NAME {
		driver = &saphana_api.SapHanaDriver{}
	} else {
		return nil, errors.New("Unsupported Database Type: " + dbType.Name)
	}

	if err := driver.Connect(dbConn, ro); err != nil {
		return nil, err
	}

	return driver, nil
}
