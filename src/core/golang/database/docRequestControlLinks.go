package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AddDocRequestControlLinkWithTx(requestId int64, controlId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO request_control_link (request_id, control_id, org_id)
		VALUES ($1, $2, $3)
	`, requestId, controlId, orgId)
	return err
}

func FindControlLinkedToDocRequest(requestId int64, orgId int32, role *core.Role) (*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	control := core.Control{}
	rows, err := dbConn.Queryx(`
		SELECT ctrl.*
		FROM process_flow_controls AS ctrl
		INNER JOIN request_control_link AS link
			ON link.control_id = ctrl.id
		WHERE link.request_id = $1 AND link.org_id = $2
	`, requestId, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	err = rows.StructScan(&control)
	return &control, err

}

func FindDocRequestsLinkedToControl(controlId int64, orgId int32, role *core.Role) ([]*core.DocumentRequest, error) {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	requests := make([]*core.DocumentRequest, 0)
	err := dbConn.Select(&requests, `
		SELECT req.*
		FROM document_requests AS req
		INNER JOIN request_control_link AS link
			ON link.request_id = req.id
		WHERE link.control_id = $1 AND link.org_id = $2
	`, controlId, orgId)
	return requests, err
}
