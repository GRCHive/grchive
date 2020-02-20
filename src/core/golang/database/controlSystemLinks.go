package database

import (
	"gitlab.com/grchive/grchive/core"
)

func FindSystemsLinkedToControl(controlId int64, orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT sys.*
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
	return systems, err
}
