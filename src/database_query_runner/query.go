package main

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/db_api"
	"gitlab.com/grchive/grchive/db_api/utility"
	"gitlab.com/grchive/grchive/webcore"
)

func runQuery(queryId int64, orgId int32) (*utility.SqlQueryResult, error) {
	query, err := database.GetSqlQueryFromId(queryId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	metadata, err := database.GetSqlMetadataFromId(query.MetadataId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	dbConn, err := database.FindDatabaseConnectionForDatabase(metadata.DbId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	if dbConn == nil {
		return nil, errors.New("Failed to find database connection.")
	}

	dbConn.Password, err = webcore.DecryptSaltedEncryptedPassword(dbConn.Password, dbConn.Salt)
	if err != nil {
		return nil, err
	}

	dbType, err := database.GetDbType(metadata.DbId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	driver, err := db_api.CreateDriver(dbType, dbConn)
	if err != nil {
		return nil, errors.New("Failed to connect to database.")
	}
	defer driver.Close()

	result, err := driver.RunQuery(query.Query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
