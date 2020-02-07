package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"time"
)

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

	return nil
}

func CreateNewDatabaseTableWithTx(schema *core.DbTable, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	return nil
}

func CreateNewDatabaseColumnWithTx(schema *core.DbColumn, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	return nil
}
