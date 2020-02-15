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
	if conn != nil {
		err = strings.Replace(err, conn.Password, "*******", -1)
	}
	database.MarkFailureRefresh(refresh.Id, err, core.ServerRole)
	return &webcore.RabbitMQError{errors.New(err), true}
}

func onRefreshSuccess(refresh *core.DbRefresh, tx *sqlx.Tx) error {
	return database.MarkSuccessfulRefreshWithTx(refresh.Id, tx, core.ServerRole)
}

func processRefreshRequest(refresh *core.DbRefresh) *webcore.RabbitMQError {
	db, err := database.GetDb(refresh.DbId, refresh.OrgId, core.ServerRole)
	if err != nil {
		return onRefreshError(nil, refresh, "Failed to get DB: "+err.Error())
	}

	dbType, err := database.GetDbType(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		return onRefreshError(nil, refresh, "Failed to get DB type: "+err.Error())
	}

	if !dbType.HasSqlSupport {
		return onRefreshError(nil, refresh, "Database type unsupported.")
	}

	conn, err := database.FindDatabaseConnectionForDatabase(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		return onRefreshError(nil, refresh, "Failed to find database connection: "+err.Error())
	}

	if conn == nil {
		return onRefreshError(nil, refresh, "Failed to find database connection.")
	}

	conn.Password, err = webcore.DecryptSaltedEncryptedPassword(conn.Password, conn.Salt)
	if err != nil {
		return onRefreshError(conn, refresh, "Failed to decrypt password: "+err.Error())
	}

	driver, err := db_api.CreateDriver(dbType, conn)
	if err != nil {
		// Don't put error here just in case there's a PW lurking around.
		return onRefreshError(conn, refresh, "Failed to connect to database.")
	}

	if !driver.ConnectionReadOnly() {
		return onRefreshError(conn, refresh, "The database user has non-read permissions.")
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

		fns, err := driver.GetFunctions(sch)
		if err != nil {
			tx.Rollback()
			return onRefreshError(conn, refresh, "Failed to get functions ["+sch.SchemaName+"]: "+err.Error())
		}

		for _, fn := range fns {
			fn.SchemaId = sch.Id
			fn.OrgId = sch.OrgId
			err = database.CreateNewDatabaseFunctionWithTx(fn, tx, core.ServerRole)
			if err != nil {
				tx.Rollback()
				return onRefreshError(conn, refresh, "Failed to store function ["+fn.Name+"]: "+err.Error())
			}
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
		return onRefreshError(conn, refresh, "Failed to commit results: "+err.Error())
	}

	return nil
}
