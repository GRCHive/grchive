package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"time"
)

func CreateNewDocumentRequestWithTx(request *core.DocumentRequest, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessManage) {
		return core.ErrorUnauthorized
	}

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

func CreateNewDocumentRequest(request *core.DocumentRequest, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := CreateNewDocumentRequestWithTx(request, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func UpdateDocumentRequest(request *core.DocumentRequest, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		UPDATE document_requests
		SET name = :name,
			description = :description
		WHERE id = :id
			AND org_id = :org_id
			AND cat_id = :cat_id
		RETURNING *
	`, request)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.StructScan(request)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func GetDocumentRequest(requestId int64, orgId int32, role *core.Role) (*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	req := core.DocumentRequest{}
	err := dbConn.Get(&req, `
		SELECT *
		FROM document_requests
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId)
	return &req, err
}

func DeleteDocumentRequest(requestId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM document_requests
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func CompleteDocumentRequest(requestId int64, orgId int32, complete bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	newTime := core.NullTime{}
	if complete {
		newTime = core.CreateNullTime(time.Now().UTC())
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		UPDATE document_requests
		SET completion_time = $3
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId, newTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetAllDocumentRequestsForVendorProduct(productId int64, orgId int32, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, `
		SELECT req.*
		FROM document_requests AS req
		INNER JOIN vendor_soc_request_link AS link
			ON req.id = link.request_id
		WHERE req.org_id = $1 AND link.vendor_product_id = $2
	`, orgId, productId)
	return requests, err
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

func FulfillDocumentRequestWithTx(requestId int64, fileId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO document_request_fulfillment (
			org_id,
			fulfilled_file_id,
			request_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
	`, orgId, fileId, requestId)
	if err != nil {
		return err
	}

	// Check if the request has any linked vendors and link if needed.
	_, err = tx.Exec(`
		INSERT INTO vendor_product_soc_reports (product_id, org_id, file_id)
		SELECT link.vendor_product_id, $2, $3
		FROM vendor_soc_request_link AS link
		INNER JOIN document_requests AS req
			ON req.id = link.request_id
		WHERE req.id = $1 AND req.org_id = $2
	`, requestId, orgId, fileId)
	if err != nil {
		return err
	}

	return nil
}

func GetFulfilledFileIdsForDocRequest(requestId int64, orgId int32, role *core.Role) ([]int64, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	ids := make([]int64, 0)
	err := dbConn.Select(&ids, `
		SELECT fulfilled_file_id
		FROM document_request_fulfillment
		WHERE org_id = $1
			AND request_id = $2
	`, orgId, requestId)
	return ids, err
}

func GetDocumentRequestComments(requestId int64, orgId int32, role *core.Role) ([]*core.Comment, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return getComments(`
		INNER JOIN document_request_comments AS drc
			ON drc.comment_id = comments.id
		WHERE drc.request_id = $1
			AND drc.org_id = $2
	`, requestId, orgId)
}

func InsertDocumentRequestComment(requestId int64, orgId int32, comment *core.Comment, role *core.Role) error {
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
		INSERT INTO document_request_comments (
			request_id,
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

func LinkRequestToVendorProductWithTx(productId int64, requestId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO vendor_soc_request_link (vendor_product_id, org_id, request_id)
		VALUES ($1, $2, $3)
	`, productId, orgId, requestId)

	return err
}
