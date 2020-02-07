package main

import (
	"flag"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/db_api"
)

func onRefreshError(refresh *core.DbRefresh, err string) {
	core.Error(err)
}

func onRefreshSuccess(refresh *core.DbRefresh, tx *sqlx.Tx) error {
	return nil
}

func main() {
	core.Init()
	database.Init()

	dbId := flag.Int64("dbId", -1, "Database ID to retrieve data for.")
	orgId := flag.Int("orgId", -1, "Organization ID to retrieve data for.")
	flag.Parse()

	if *dbId == -1 {
		core.Error("No database ID specified.")
	}

	if *orgId == -1 {
		core.Error("No organization ID specified.")
	}

	db, err := database.GetDb(*dbId, int32(*orgId), core.ServerRole)
	if err != nil {
		core.Error("Failed to Get DB: " + err.Error())
	}

	dbType, err := database.GetDbType(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		core.Error("Failed to Get DB Type: " + err.Error())
	}

	if !dbType.HasSqlSupport {
		core.Error("Database does not have SQL API support.")
	}

	conn, err := database.FindDatabaseConnectionForDatabase(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		core.Error("Failed to Get DB Connection: " + err.Error())
	}

	driver, err := db_api.CreateDriver(dbType, conn)
	if err != nil {
		core.Error("Failed to connect to database: " + err.Error())
	}

	if !driver.ConnectionReadOnly() {
		core.Error("The databaser user provided has non-read permissions.")
	}

	tx := database.CreateTx()

	// Create a DbRefresh object to track this fetch operation.
	refresh, err := database.CreateNewDatabaseRefresh(db.Id, db.OrgId, core.ServerRole)
	if err != nil {
		core.Error("Failed to create DB Refresh: " + err.Error())
	}

	schemas, err := driver.GetSchemas()
	if err != nil {
		tx.Rollback()
		onRefreshError(refresh, "Failed to get schemas: "+err.Error())
	}

	for _, sch := range schemas {
		sch.RefreshId = refresh.Id
		err = database.CreateNewDatabaseSchemaWithTx(sch, tx, core.ServerRole)
		if err != nil {
			tx.Rollback()
			onRefreshError(refresh, "Failed to store schema ["+sch.SchemaName+"]: "+err.Error())
		}

		tables, err := driver.GetTables(sch)
		if err != nil {
			tx.Rollback()
			onRefreshError(refresh, "Failed to get tables ["+sch.SchemaName+"]: "+err.Error())
		}

		for _, tbl := range tables {
			tbl.SchemaId = sch.Id
			err = database.CreateNewDatabaseTableWithTx(tbl, tx, core.ServerRole)
			if err != nil {
				tx.Rollback()
				onRefreshError(refresh, "Failed to store table ["+tbl.TableName+"]: "+err.Error())
			}

			columns, err := driver.GetColumns(tbl)
			if err != nil {
				tx.Rollback()
				onRefreshError(refresh, "Failed to get columns ["+tbl.TableName+"]: "+err.Error())
			}

			for _, col := range columns {
				col.TableId = tbl.Id
				err = database.CreateNewDatabaseColumnWithTx(col, tx, core.ServerRole)
				if err != nil {
					tx.Rollback()
					onRefreshError(refresh, "Failed to store column ["+col.ColumnName+"]: "+err.Error())
				}
			}
		}
	}

	err = onRefreshSuccess(refresh, tx)
	if err != nil {
		onRefreshError(refresh, "Failed to mark successful refresh: "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		core.Error("Failed to commit: " + err.Error())
	}
}
