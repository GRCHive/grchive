package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/db_api"
	"gitlab.com/grchive/grchive/webcore"
	"strings"
	"time"
)

func onRefreshError(conn *core.DatabaseConnection, refresh *core.DbRefresh, err string) *webcore.RabbitMQError {
	if conn != nil && conn.Password != "" {
		err = strings.Replace(err, conn.Password, "*******", -1)
	}
	database.MarkFailureRefresh(refresh.Id, err, core.ServerRole)
	return &webcore.RabbitMQError{errors.New(err), false}
}

func onRefreshSuccess(refresh *core.DbRefresh, hasDiff bool, tx *sqlx.Tx) error {
	return database.MarkSuccessfulRefreshWithTx(refresh.Id, hasDiff, tx, core.ServerRole)
}

func processRefreshRequest(refresh *core.DbRefresh, sendEmail bool) *webcore.RabbitMQError {
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

	driver, err := db_api.CreateDriver(dbType, conn, true)
	if err != nil {
		return onRefreshError(conn, refresh, "Failed to connect to database: "+err.Error())
	}
	defer driver.Close()

	latestRefresh, err := database.GetLatestCompleteDatabaseRefresh(refresh.DbId, refresh.OrgId)
	if err != nil {
		return onRefreshError(conn, refresh, "Failed to get latest refresh: "+err.Error())
	}

	var latestDbState *core.FullDbState
	if latestRefresh != nil {
		latestDbState, err = webcore.GetDbStateFromRefresh(latestRefresh.Id, latestRefresh.OrgId)
		if err != nil {
			return onRefreshError(conn, refresh, "Failed to get latest db state: "+err.Error())
		}
	}

	err = driver.ConstructState()
	if err != nil {
		return onRefreshError(conn, refresh, "Failed to construct state: "+err.Error())
	}

	currentDbState := driver.GetState()
	tx := database.CreateTx()

	startCommit := time.Now()
	hasDiff := true
	err = database.WrapTx(tx, func() error {
		schemas := currentDbState.AllSchemas()
		for _, sch := range schemas {
			sch.DbSchema.OrgId = refresh.OrgId
			sch.DbSchema.RefreshId = refresh.Id

			err = database.CreateNewDatabaseSchemaWithTx(&sch.DbSchema, tx, core.ServerRole)
			if err != nil {
				return err
			}

			fns := sch.AllFunctions()
			for _, fn := range fns {
				fn.OrgId = refresh.OrgId
				fn.SchemaId = sch.DbSchema.Id
				err = database.CreateNewDatabaseFunctionWithTx(fn, tx, core.ServerRole)
				if err != nil {
					return err
				}
			}

			tables := sch.AllTables()
			for _, tbl := range tables {
				tbl.OrgId = refresh.OrgId
				tbl.SchemaId = sch.DbSchema.Id
				err = database.CreateNewDatabaseTableWithTx(tbl, tx, core.ServerRole)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}, func() error {
		if latestDbState != nil {
			hasDiff = currentDbState.HasDiff(latestDbState)
		}

		return onRefreshSuccess(refresh, hasDiff, tx)
	})

	if err != nil {
		return onRefreshError(conn, refresh, "Failed to perform refresh: "+err.Error())
	}

	core.Debug(fmt.Sprintf(
		"\tFinish Commit: %f seconds",
		float64(time.Since(startCommit).Milliseconds())/1000.0,
	))

	// Ideally this would happen elsewhere and be sent out as a general "event"
	// but our event/notification system is kind of rigid at the moment and needs
	// a revamp. We can ignore errors here.
	if hasDiff && sendEmail {
		err = sendDbSchemaChangeEmails(db)
		if err != nil {
			core.Warning("Failed to send Db schema change emails: " + err.Error())
		}
	}

	return nil
}
