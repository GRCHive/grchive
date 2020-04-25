package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func CreateGenericRequestWithTx(tx *sqlx.Tx, req *core.GenericRequest) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO generic_requests (org_id, upload_time, upload_user_id, name, assignee, due_date, description)
		VALUES (:org_id, :upload_time, :upload_user_id, :name, :assignee, :due_date, :description)
		RETURNING id
	`, req)

	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	return rows.Scan(&req.Id)
}

func EditGenericRequest(reqId int64, orgId int32, newData core.GenericRequest, role *core.Role) error {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	return WrapTx(tx, func() error {
		_, err := tx.Exec(`
			UPDATE generic_requests
				SET name = $3,
					description = $4,
					assignee = $5,
					due_date = $6
			WHERE id = $1
				AND org_id = $2
		`, reqId, orgId, newData.Name, newData.Description, newData.Assignee, newData.DueDate)
		return err
	})
}

func DeleteGenericRequest(reqId int64, orgId int32, role *core.Role) error {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	return WrapTx(tx, func() error {
		_, err := tx.Exec(`
			DELETE FROM generic_requests
			WHERE id = $1
				AND org_id = $2
		`, reqId, orgId)
		return err
	})
}

func LinkScheduledTaskToRequestWithTx(tx *sqlx.Tx, taskId int64, requestId int64) error {
	_, err := tx.Exec(`
		INSERT INTO request_to_scheduled_task_link (request_id, task_id)
		VALUES ($2, $1)
	`, taskId, requestId)
	return err
}

func LinkScriptRunToRequestWithTx(tx *sqlx.Tx, runId int64, requestId int64) error {
	_, err := tx.Exec(`
		INSERT INTO request_to_script_run_link (request_id, run_id)
		VALUES ($2, $1)
	`, runId, requestId)
	return err
}

func GetGenericRequestsForScriptsInOrg(orgId int32, role *core.Role) ([]*core.GenericRequest, error) {
	reqs := make([]*core.GenericRequest, 0)
	err := dbConn.Select(&reqs, `
		SELECT DISTINCT(req.*)
		FROM generic_requests AS req
		LEFT JOIN request_to_scheduled_task_link AS tl
			ON tl.request_id = req.id
		LEFT JOIN scheduled_task_script_links AS tsl
			ON tsl.event_id = tl.task_id
		LEFT JOIN request_to_script_run_link AS rl
			ON rl.request_id = req.id
		WHERE req.org_id = $1
			AND (rl.run_id IS NOT NULL OR tsl.link_id IS NOT NULL)
		ORDER BY req.upload_time DESC
	`, orgId)
	return reqs, err
}

func GetGenericRequestFromId(reqId int64) (*core.GenericRequest, error) {
	req := core.GenericRequest{}
	err := dbConn.Get(&req, `
		SELECT *
		FROM generic_requests
		WHERE id = $1
	`, reqId)
	return &req, err
}

func GetCodeFromScriptRequestId(reqId int64, role *core.Role) (*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.ManagedCode{}
	err := dbConn.Get(&code, `
		SELECT DISTINCT(code.*)
		FROM managed_code AS code
		INNER JOIN code_to_client_scripts_link AS link
			ON link.code_id = code.id
		LEFT JOIN script_runs AS run
			ON run.link_id = link.id
		LEFT JOIN request_to_script_run_link AS rl
			ON rl.run_id = run.id
		LEFT JOIN scheduled_task_script_links AS tsl
			ON tsl.link_id = link.id
		LEFT JOIN request_to_scheduled_task_link AS sl
			ON sl.task_id = tsl.event_id
		LEFT JOIN generic_requests AS req
			ON (req.id = rl.request_id OR req.id = sl.request_id)
		WHERE req.id = $1
			AND (rl.request_id IS NOT NULL OR
				sl.request_id IS NOT NULL)
	`, reqId)
	return &code, err
}

func InsertGenericApproval(approval *core.GenericApproval, role *core.Role) error {
	rows, err := dbConn.NamedQuery(`
		INSERT INTO generic_approval (request_id, response_time, responder_user_id, response, reason)
		VALUES (:request_id, :response_time, :responder_user_id, :response, :reason)
		RETURNING id
	`, approval)

	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	return rows.Scan(&approval.Id)
}

func GetGenericApprovalForRequest(requestId int64, orgId int32, role *core.Role) (*core.GenericApproval, error) {
	rows, err := dbConn.Queryx(`
		SELECT appr.*
		FROM generic_approval AS appr
		INNER JOIN generic_requests AS req
			ON req.id = appr.request_id
		WHERE appr.request_id = $1 
			AND req.org_id = $2
	`, requestId, orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	approval := core.GenericApproval{}
	err = rows.StructScan(&approval)
	return &approval, err
}

func GetGenericRequestCommentThreadId(requestId int64, role *core.Role) (int64, error) {
	threadId := int64(0)
	err := dbConn.Get(&threadId, `
		SELECT thread_id
		FROM generic_request_comment_threads
		WHERE request_id = $1
	`, requestId)
	if err != nil {
		return -1, err
	}
	return threadId, nil
}

func InsertGenericRequestComment(requestId int64, comment *core.Comment, role *core.Role) error {
	threadId, err := GetGenericRequestCommentThreadId(requestId, role)
	if err != nil {
		return err
	}

	tx := dbConn.MustBegin()
	err = insertCommentWithTx(comment, threadId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetGenericRequestComments(reqId int64, orgId int32, role *core.Role) ([]*core.Comment, error) {
	return getComments(`
		INNER JOIN generic_request_comment_threads AS rct
			ON rct.thread_id = t.id
		WHERE rct.request_id = $1
			AND t.org_id = $2
	`, reqId, orgId)
}
