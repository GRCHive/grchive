package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AddDocRequestDocCatLinkWithTx(requestId int64, catId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO request_doc_cat_link (request_id, cat_id, org_id)
		VALUES ($1, $2, $3)
	`, requestId, catId, orgId)
	return err
}

func FindDocCatLinkedToDocRequest(requestId int64, orgId int32, role *core.Role) (*core.ControlDocumentationCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	cat := core.ControlDocumentationCategory{}
	rows, err := dbConn.Queryx(`
		SELECT cat.*
		FROM process_flow_control_documentation_categories AS cat
		INNER JOIN request_doc_cat_link AS link
			ON link.cat_id = cat.id
		WHERE link.request_id = $1 AND link.org_id = $2
	`, requestId, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	err = rows.StructScan(&cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func FindDocRequestsLinkedToDocCat(catId int64, orgId int32, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, `
		SELECT req.*
		FROM document_requests AS req
		INNER JOIN request_doc_cat_link AS link
			ON link.request_id = req.id
		WHERE link.cat_id = $1 AND link.org_id = $2
	`, catId, orgId)

	if err != nil {
		return nil, err
	}

	return requests, nil
}
