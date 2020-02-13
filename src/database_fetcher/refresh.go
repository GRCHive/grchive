package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/db_api"
	"gitlab.com/grchive/grchive/webcore"
	"strings"
)

func onRefreshError(conn *core.DatabaseConnection, refresh *core.DbRefresh, err string) *webcore.RabbitMQError {
	err = strings.Replace(err, conn.Password, "*******", -1)
	database.MarkFailureRefresh(refresh.Id, err, core.ServerRole)
	return &webcore.RabbitMQError{errors.New(err), true}
}

func onRefreshSuccess(refresh *core.DbRefresh, tx *sqlx.Tx) error {
	return database.MarkSuccessfulRefreshWithTx(refresh.Id, tx, core.ServerRole)
}

func processRefreshRequest(refresh *core.DbRefresh) *webcore.RabbitMQError {
	db, err := database.GetDb(refresh.DbId, refresh.OrgId, core.ServerRole)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	dbType, err := database.GetDbType(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	if !dbType.HasSqlSupport {
		return &webcore.RabbitMQError{errors.New("Database type unsupported."), false}
	}

	conn, err := database.FindDatabaseConnectionForDatabase(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	if conn == nil {
		return &webcore.RabbitMQError{errors.New("Failed to find database connection."), false}
	}

	conn.Password, err = webcore.DecryptSaltedEncryptedPassword(conn.Password, conn.Salt)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	driver, err := db_api.CreateDriver(dbType, conn)
	if err != nil {
		// Don't put error here just in case there's a PW lurking around.
		return &webcore.RabbitMQError{errors.New("Failed to connect to database."), false}
	}

	if !driver.ConnectionReadOnly() {
		return &webcore.RabbitMQError{errors.New("The database user has non-read permissions."), false}
	}

	tx := database.CreateTx()

	schemas, err := driver.GetSchemas()
	if err != nil {
		tx.Rollback()
		return onRefreshError(conn, refresh, "Failed to get schemas: "+err.Error())
	}

	for _, sch := range schemas {
		sch.RefreshId = refresh.Id
		sch.OrgId = refresh.OrgId
		err = database.CreateNewDatabaseSchemaWithTx(sch, tx, core.ServerRole)
		if err != nil {
			tx.Rollback()
			return onRefreshError(conn, refresh, "Failed to store schema ["+sch.SchemaName+"]: "+err.Error())
		}

		tables, err := driver.GetTables(sch)
		if err != nil {
			tx.Rollback()
			return onRefreshError(conn, refresh, "Failed to get tables ["+sch.SchemaName+"]: "+err.Error())
		}

		for _, tbl := range tables {
			tbl.SchemaId = sch.Id
			tbl.OrgId = sch.OrgId
			err = database.CreateNewDatabaseTableWithTx(tbl, tx, core.ServerRole)
			if err != nil {
				tx.Rollback()
				return onRefreshError(conn, refresh, "Failed to store table ["+tbl.TableName+"]: "+err.Error())
			}

			columns, err := driver.GetColumns(sch, tbl)
			if err != nil {
				tx.Rollback()
				return onRefreshError(conn, refresh, "Failed to get columns ["+tbl.TableName+"]: "+err.Error())
			}

			for _, col := range columns {
				col.TableId = tbl.Id
				col.OrgId = tbl.OrgId
				err = database.CreateNewDatabaseColumnWithTx(col, tx, core.ServerRole)
				if err != nil {
					tx.Rollback()
					return onRefreshError(conn, refresh, "Failed to store column ["+col.ColumnName+"]: "+err.Error())
				}
			}
		}
	}

	err = onRefreshSuccess(refresh, tx)
	if err != nil {
		tx.Rollback()
		return onRefreshError(conn, refresh, "Failed to mark successful refresh: "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	return nil
}
