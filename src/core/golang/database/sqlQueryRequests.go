package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func CreateNewSqlQueryRequestWithTx(request *core.DbSqlQueryRequest, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessManage) {
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

func GetAllSqlRequestsForDb(dbId int64, orgId int32, role *core.Role) ([]*core.DbSqlQueryRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbSqlQueryRequest, 0)
	err := dbConn.Select(&data, `
		SELECT req.*
		FROM database_sql_query_requests AS req
		INNER JOIN database_sql_queries AS query
			ON req.query_id = query.id
		INNER JOIN database_sql_metadata AS meta
			ON query.metadata_id = meta.id
		WHERE meta.db_id = $1 AND meta.org_id = $2
		ORDER BY req.upload_time DESC
	`, dbId, orgId)
	return data, err
}

func GetAllSqlRequestsForOrg(orgId int32, role *core.Role) ([]*core.DbSqlQueryRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbSqlQueryRequest, 0)
	err := dbConn.Select(&data, `
		SELECT req.*
		FROM database_sql_query_requests AS req
		WHERE req.org_id = $1
		ORDER BY req.upload_time DESC
	`, orgId)
	return data, err
}

func GetSqlRequestStatus(requestId int64, orgId int32, role *core.Role) (*core.DbSqlQueryRequestApproval, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT *
		FROM database_sql_query_requests_approvals
		WHERE request_id = $1 AND org_id = $2
	`, requestId, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	approval := core.DbSqlQueryRequestApproval{}
	err = rows.StructScan(&approval)
	if err != nil {
		return nil, err
	}

	return &approval, nil
}
