package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"time"
)

func GetAllDatabaseSchemaForRefresh(refreshId int64, orgId int32, role *core.Role) ([]*core.DbSchema, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbSchema, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_schemas
		WHERE refresh_id = $1 AND org_id = $2
		ORDER BY schema_name ASC
	`, refreshId, orgId)
	return data, err
}

func GetDatabaseRefresh(refreshId int64, orgId int32, role *core.Role) (*core.DbRefresh, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	refresh := core.DbRefresh{}
	err := dbConn.Get(&refresh, `
		SELECT *
		FROM database_refresh
		WHERE id = $1 AND org_id = $2
	`, refreshId, orgId)
	return &refresh, err
}

func GetLatestCompleteDatabaseRefresh(dbId int64, orgId int32) (*core.DbRefresh, error) {
	rows, err := dbConn.Queryx(`
		SELECT *
		FROM database_refresh
		WHERE refresh_success = true
			AND db_id = $1
			AND org_id = $2
		ORDER BY id DESC
		LIMIT 1
	`, dbId, orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	refresh := core.DbRefresh{}
	err = rows.StructScan(&refresh)
	return &refresh, err
}

func DeleteDatabaseRefresh(refreshId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM database_refresh
		WHERE id = $1 AND org_id = $2
	`, refreshId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func MarkSuccessfulRefreshWithTx(refreshId int64, hasDiff bool, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		UPDATE database_refresh
		SET refresh_success = true,
			refresh_finish_time = $2,
			refresh_has_diff = $3
		WHERE id = $1
	`, refreshId, time.Now().UTC(), hasDiff)
	return err
}

func MarkFailureRefresh(refreshId int64, failureReason string, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	_, err := tx.Exec(`
		UPDATE database_refresh
		SET refresh_success = false,
			refresh_errors = $2,
			refresh_finish_time = $3
		WHERE id = $1
	`, refreshId, failureReason, time.Now().UTC())

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func CreateNewDatabaseRefreshWithTx(dbId int64, orgId int32, tx *sqlx.Tx, role *core.Role) (*core.DbRefresh, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return nil, core.ErrorUnauthorized
	}

	refresh := core.DbRefresh{
		DbId:           dbId,
		OrgId:          orgId,
		RefreshTime:    core.CreateNullTime(time.Now().UTC()),
		RefreshSuccess: false,
		RefreshErrors:  "N/A",
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_refresh (db_id, org_id, refresh_time, refresh_success, refresh_errors)
		VALUES (:db_id, :org_id, :refresh_time, :refresh_success, :refresh_errors)
		RETURNING id
	`, refresh)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&refresh.Id)
	return &refresh, err
}

func CreateNewDatabaseRefresh(dbId int64, orgId int32, role *core.Role) (*core.DbRefresh, error) {
	tx := dbConn.MustBegin()
	refresh, err := CreateNewDatabaseRefreshWithTx(dbId, orgId, tx, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return refresh, tx.Commit()
}

func CreateNewDatabaseSchemaWithTx(schema *core.DbSchema, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_schemas (org_id, refresh_id, schema_name)
		VALUES (:org_id, :refresh_id, :schema_name)
		RETURNING id
	`, schema)

	if err != nil {
		return err
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&schema.Id)
	return err
}

func CreateNewDatabaseTableWithTx(table *core.DbTable, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_tables (org_id, schema_id, table_name)
		VALUES (:org_id, :schema_id, :table_name)
		RETURNING id
	`, table)

	if err != nil {
		return err
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&table.Id)
	return err
}

func CreateNewDatabaseColumnWithTx(column *core.DbColumn, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_columns (org_id, table_id, column_name, column_type)
		VALUES (:org_id, :table_id, :column_name, :column_type)
		RETURNING id
	`, column)

	if err != nil {
		return err
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&column.Id)
	return err
}

func CreateNewDatabaseFunctionWithTx(fn *core.DbFunction, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_functions (org_id, schema_id, name, src, ret_type)
		VALUES (:org_id, :schema_id, :name, :src, :ret_type)
		RETURNING id
	`, fn)

	if err != nil {
		return err
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&fn.Id)
	return err
}

func GetAllDatabaseRefresh(dbId int64, orgId int32, role *core.Role) ([]*core.DbRefresh, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbRefresh, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_refresh
		WHERE db_id = $1 AND org_id = $2
		ORDER BY refresh_time DESC
	`, dbId, orgId)
	return data, err
}
