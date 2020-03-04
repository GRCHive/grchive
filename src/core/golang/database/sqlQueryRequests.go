package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"time"
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
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = CreateNewSqlQueryRequestWithTx(request, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func UpdateSqlQueryRequestWithTx(request *core.DbSqlQueryRequest, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		UPDATE database_sql_query_requests
		SET name = :name,
			description = :description
		WHERE id = :id
			AND org_id = :org_id
		RETURNING *
	`, request)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(request)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSqlQueryRequest(request *core.DbSqlQueryRequest, role *core.Role) error {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = UpdateSqlQueryRequestWithTx(request, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteSqlQueryRequest(requestId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM database_sql_query_requests
		WHERE id = $1 AND org_id = $2
	`, requestId, orgId)
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

func GetSqlRequest(requestId int64, orgId int32, role *core.Role) (*core.DbSqlQueryRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	req := core.DbSqlQueryRequest{}
	err := dbConn.Get(&req, `
		SELECT *
		FROM database_sql_query_requests
		WHERE id = $1 AND org_id = $2
	`, requestId, orgId)
	return &req, err
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

func UpdateRequestStatusWithTx(approval *core.DbSqlQueryRequestApproval, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_sql_query_requests_approvals (request_id, org_id, response_time, responder_user_id, response, reason)
		VALUES (:request_id, :org_id, :response_time, :responder_user_id, :response, :reason)
		RETURNING *
	`, approval)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()

	err = rows.StructScan(approval)
	if err != nil {
		return err
	}
	return nil
}

func CreateNewRunCodeWithTx(runCode *core.DbSqlQueryRunCode, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.NamedExec(`
		INSERT INTO database_sql_query_run_codes (request_id, org_id, expiration_time, used_time, hashed_code, salt)
		VALUES (:request_id, :org_id, :expiration_time, NULL, :hashed_code, :salt)
		ON CONFLICT (request_id) DO UPDATE
			SET expiration_time = EXCLUDED.expiration_time,
				used_time = EXCLUDED.used_time,
				hashed_code = EXCLUDED.hashed_code,
				salt = EXCLUDED.salt
	`, runCode)

	return err
}

func FindRunCodeForQueryForUser(queryId int64, orgId int32, userId int64, role *core.Role) (*core.DbSqlQueryRunCode, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.DbSqlQueryRunCode{}
	err := dbConn.Get(&code, `
		SELECT code.*
		FROM database_sql_query_run_codes AS code
		INNER JOIN database_sql_query_requests AS req
			ON code.request_id = req.id
		WHERE req.query_id = $1
			AND req.org_id = $2
			AND req.upload_user_id = $3
			AND code.used_time IS NULL
	`, queryId, orgId, userId)
	return &code, err
}

func MarkRunCodeAsUsed(hashCode string, requestId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlRequest, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		UPDATE database_sql_query_run_codes
		SET used_time = $4
		WHERE hashed_code = $1 AND request_id = $2 AND org_id = $3
	`, hashCode, requestId, orgId, time.Now().UTC())

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetSqlRequestComments(requestId int64, orgId int32, role *core.Role) ([]*core.Comment, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return getComments(`
		INNER JOIN sql_request_comments AS src
			ON src.comment_id = comments.id
		WHERE src.sql_request_id = $1
			AND src.org_id = $2
	`, requestId, orgId)
}

func InsertSqlRequestComment(requestId int64, orgId int32, comment *core.Comment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	err := insertCommentWithTx(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO sql_request_comments (
			sql_request_id,
			org_id,
			comment_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
	`, requestId, orgId, comment.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
