package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func CreateNewDocumentRequest(request *core.DocumentRequest, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO document_requests(
			name,
			description,
			cat_id,
			org_id,
			requested_user_id,
			request_time
		)
		VALUES (
			:name,
			:description,
			:cat_id,
			:org_id,
			:requested_user_id,
			:request_time
		)
		RETURNING id
	`, request)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&request.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func GetDocumentRequest(requestId int64, catId int64, orgId int32, role *core.Role) (*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	req := core.DocumentRequest{}
	err := dbConn.Get(&req, `
		SELECT *
		FROM document_requests
		WHERE id = $1
			AND cat_id = $2
			AND org_id = $3
	`, requestId, catId, orgId)
	return &req, err
}

func GetAllDocumentRequestsForOrganization(orgId int32, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, `
		SELECT *
		FROM document_requests
		WHERE org_id = $1
	`, orgId)
	return requests, err
}

func GetAllDocumentRequestsForDocCat(catId int64, orgId int32, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, `
		SELECT *
		FROM document_requests
		WHERE cat_id = $1
			AND org_id = $2
	`, catId, orgId)
	return requests, err
}

func FulfillDocumentRequest(fileId int64, catId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO document_request_fulfillment (
			cat_id,
			org_id,
			fulfilled_file_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
	`, catId, orgId, fileId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func AddDocRequestComment(comment *core.DocumentRequestComment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO document_request_comments (
			user_id,
			text,
			post_time,
			cat_id,
			org_id,
			request_id
		)
		VALUES (
			:user_id,
			:text,
			:post_time,
			:cat_id,
			:org_id,
			:request_id
		)
		RETURNING id
	`, comment)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&comment.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func GetAllDocumentRequestComments(requestId int64, orgId int32, role *core.Role) ([]*core.DocumentRequestComment, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	comments := make([]*core.DocumentRequestComment, 0)
	err := dbConn.Select(&comments, `
		SELECT *
		FROM document_request_comments
		WHERE org_id = $1
			AND request_id = $2
	`, orgId, requestId)
	return comments, err
}