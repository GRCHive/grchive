package database

import (
	"gitlab.com/grchive/grchive/core"
)

func FindSystemsLinkedToRisk(riskId int64, orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) ||
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
		INNER JOIN process_flow_risks AS risk
			ON risk.id = rn.risk_id
		WHERE risk.id = $1 AND risk.org_id = $2
	`, riskId, orgId)
	return systems, err
}
