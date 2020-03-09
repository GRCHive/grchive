package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
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

func FindRisksLinkedToSystem(systemId int64, orgId int32, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	risks := make([]*core.Risk, 0)
	err := dbConn.Select(&risks, `
		SELECT DISTINCT risk.*
		FROM process_flow_risks AS risk
		INNER JOIN process_flow_risk_node AS rn
			ON rn.risk_id = risk.id
		INNER JOIN node_system_link AS nsl
			ON nsl.node_id = rn.node_id
		INNER JOIN systems AS sys
			ON sys.id = nsl.system_id
		WHERE sys.id = $1
			AND sys.org_id = $2
			AND risk.org_id = $2
	`, systemId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, r := range risks {
		err = LogAuditSelectWithTx(orgId, core.ResourceRisk, strconv.FormatInt(r.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return risks, tx.Commit()
}
