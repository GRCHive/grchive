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
		SELECT DISTINCT sys.*
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

func FindRisksLinkedToSystem(systemId int64, orgId int32, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	risks := make([]*core.Risk, 0)
	err := dbConn.Select(&risks, `
		SELECT DISTINCT
			risk.id,
			risk.name,
			risk.description,
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name"
		FROM process_flow_risks AS risk
		INNER JOIN process_flow_risk_node AS rn
			ON rn.risk_id = risk.id
		INNER JOIN node_system_link AS nsl
			ON nsl.node_id = rn.node_id
		INNER JOIN systems AS sys
			ON sys.id = nsl.system_id
		INNER JOIN organizations AS org
			ON sys.org_id = org.id
		WHERE sys.id = $1
			AND sys.org_id = $2
			AND risk.org_id = $2
	`, systemId, orgId)
	return risks, err
}
