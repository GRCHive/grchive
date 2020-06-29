package database

import (
	"fmt"
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
			org_id,
			requested_user_id,
			request_time,
			assignee,
			due_date
		)
		VALUES (
			:name,
			:description,
			:org_id,
			:requested_user_id,
			:request_time,
			:assignee,
			:due_date
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
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = CreateNewDocumentRequestWithTx(request, role, tx)
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	rows, err := tx.NamedQuery(`
		UPDATE document_requests
		SET name = :name,
			description = :description,
			assignee = :assignee,
			due_date = :due_date
		WHERE id = :id
			AND org_id = :org_id
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
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func DeleteDocumentRequest(requestId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
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

func MarkDocumentRequestProgressWithTx(requestId int64, orgId int32, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE document_requests
		SET progress_time = NOW()
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId)
	return err
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
		ORDER BY req.id DESC
	`, orgId, productId)

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func GetAllDocumentRequestsForOrganization(orgId int32, filter core.DocRequestFilterData, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, fmt.Sprintf(`
		SELECT *
		FROM document_requests
		WHERE org_id = $1
			AND %s
		ORDER BY id DESC
	`, buildDocRequestFilter("document_requests", filter)), orgId)

	if err != nil {
		return nil, err
	}

	return requests, nil
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
		ON CONFLICT (org_id, request_id, fulfilled_file_id) DO NOTHING
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
		ON CONFLICT (product_id, file_id) DO NOTHING
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
		INNER JOIN document_request_comment_threads AS drc
			ON drc.thread_id = t.id
		WHERE drc.request_id = $1
			AND drc.org_id = $2
	`, requestId, orgId)
}

func GetDocumentRequestCommentThreadId(requestId int64, orgId int32, role *core.Role) (int64, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) {
		return -1, core.ErrorUnauthorized
	}

	threadId := int64(-1)
	err := dbConn.Get(&threadId, `
		SELECT thread_id
		FROM document_request_comment_threads
		WHERE request_id = $1 AND org_id = $2
	`, requestId, orgId)
	return threadId, err
}

func InsertDocumentRequestComment(requestId int64, orgId int32, comment *core.Comment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	threadId := int64(0)
	err := dbConn.Get(&threadId, `
		SELECT thread_id
		FROM document_request_comment_threads
		WHERE request_id = $1 AND org_id = $2
	`, requestId, orgId)
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

func GetVendorProductIdForDocRequest(requestId int64, orgId int32) (int64, int64, error) {
	rows, err := dbConn.Queryx(`
		SELECT prod.vendor_id, link.vendor_product_id
		FROM vendor_soc_request_link AS link
		INNER JOIN vendor_products AS prod
			ON prod.id = link.vendor_product_id
		WHERE link.request_id = $1 AND link.org_id = $2
	`, requestId, orgId)
	if err != nil {
		return -1, -1, err
	}

	defer rows.Close()
	if !rows.Next() {
		return -1, -1, nil
	}
	vendorId := int64(0)
	productId := int64(0)
	err = rows.Scan(&vendorId, &productId)
	return vendorId, productId, err
}

func CompleteDocumentRequestWithTx(tx *sqlx.Tx, requestId int64, orgId int32) error {
	tm := time.Now().UTC()

	_, err := tx.Exec(`
		UPDATE document_requests
		SET completion_time = $3
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId, tm)

	return err
}

func ReopenDocumentRequestWithTx(tx *sqlx.Tx, requestId int64, orgId int32) error {
	tm := time.Now().UTC()

	_, err := tx.Exec(`
		UPDATE document_requests
		SET feedback_time = $3,
			approve_time = NULL,
			approve_user_id = NULL
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId, tm)

	return err
}

func ApproveDocumentRequestWithTx(tx *sqlx.Tx, requestId int64, orgId int32, approveUserId int64) error {
	tm := time.Now().UTC()

	_, err := tx.Exec(`
		UPDATE document_requests
		SET approve_time = $3,
			approve_user_id = $4
		WHERE id = $1
			AND org_id = $2
	`, requestId, orgId, tm, approveUserId)

	return err
}
