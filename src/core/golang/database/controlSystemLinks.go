package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func FindSystemsLinkedToControl(controlId int64, orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT DISTINCT sys.*
		FROM systems AS sys
		INNER JOIN node_system_link AS nsl
			ON nsl.system_id = sys.id
		INNER JOIN process_flow_risk_node AS rn
			ON rn.node_id = nsl.node_id
		INNER JOIN process_flow_risk_control AS rc
			ON rc.risk_id = rn.risk_id
		INNER JOIN process_flow_controls AS control 
			ON control.id = rc.control_id
		WHERE control.id = $1 AND control.org_id = $2
	`, controlId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, s := range systems {
		err = LogAuditSelectWithTx(orgId, core.ResourceSystem, strconv.FormatInt(s.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return systems, tx.Commit()
}

func FindControlsLinkedToSystem(systemId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	controls := make([]*core.Control, 0)
	err := dbConn.Select(&controls, `
		SELECT DISTINCT control.*
		FROM process_flow_controls AS control
		INNER JOIN process_flow_risk_control AS rc
			ON rc.control_id = control.id
		INNER JOIN process_flow_risk_node AS rn
			ON rn.risk_id = rc.risk_id
		INNER JOIN node_system_link AS nsl
			ON nsl.node_id = rn.node_id
		INNER JOIN systems AS sys
			ON sys.id = nsl.system_id
		WHERE sys.id = $1 AND sys.org_id = $2
	`, systemId, orgId)

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, c := range controls {
		err = LogAuditSelectWithTx(orgId, core.ResourceControl, strconv.FormatInt(c.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return controls, tx.Commit()
}
