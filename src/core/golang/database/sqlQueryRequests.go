package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func CreateNewSqlQueryRequestWithTx(request *core.DbSqlQueryRequest, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_sql_query_requests (query_id, upload_time, upload_user_id, org_id, name, description)
		VALUES (:query_id, :upload_time, :upload_user_id, :org_id, :name, :description)
		RETURNING id
	`, request)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&request.Id)
	if err != nil {
		return err
	}
	return nil
}

func CreateNewSqlQueryRequest(request *core.DbSqlQueryRequest, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := CreateNewSqlQueryRequestWithTx(request, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
