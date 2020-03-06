package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func FindNodeRiskRelationships(flowId int64, role *core.Role) ([]*core.NodeRiskRelationship, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	relationships := make([]*core.NodeRiskRelationship, 0)
	err := dbConn.Select(&relationships, `
		SELECT DISTINCT risknode.*
		FROM process_flow_risk_node AS risknode
		INNER JOIN process_flow_nodes AS node
			ON risknode.node_id = risknode.node_id
		INNER JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		WHERE flow.id = $1
	`, flowId)
	return relationships, err
}

func FindFlowsRelatedToRisk(riskId int64, role *core.Role) ([]*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	flows := make([]*core.ProcessFlow, 0)
	err := dbConn.Select(&flows, `
		SELECT DISTINCT
			flow.id,
			flow.name,
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name",
			flow.description,
			flow.created_time,
			flow.last_updated_time
		FROM process_flows AS flow
		INNER JOIN process_flow_nodes AS node
			ON node.process_flow_id = flow.id
		INNER JOIN process_flow_risk_node AS rn
			ON rn.node_id = node.id
		INNER JOIN organizations AS org
			ON flow.org_id = org.id
		WHERE rn.risk_id = $1
	`, riskId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, f := range flows {
		err = LogAuditSelectWithTx(f.Org.Id, core.ResourceProcessFlow, strconv.FormatInt(f.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return flows, tx.Commit()
}

func FindFlowsRelatedToControl(controlId int64, role *core.Role) ([]*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	flows := make([]*core.ProcessFlow, 0)
	err := dbConn.Select(&flows, `
		SELECT DISTINCT
			flow.id,
			flow.name,
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name",
			flow.description,
			flow.created_time,
			flow.last_updated_time
		FROM process_flows AS flow
		INNER JOIN process_flow_nodes AS node
			ON node.process_flow_id = flow.id
		INNER JOIN process_flow_control_node AS cn
			ON cn.node_id = node.id
		INNER JOIN organizations AS org
			ON flow.org_id = org.id
		WHERE cn.control_id = $1
	`, controlId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, f := range flows {
		err = LogAuditSelectWithTx(f.Org.Id, core.ResourceProcessFlow, strconv.FormatInt(f.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return flows, tx.Commit()
}

func FindNodeControlRelationships(flowId int64, role *core.Role) ([]*core.NodeControlRelationship, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	relationships := make([]*core.NodeControlRelationship, 0)
	err := dbConn.Select(&relationships, `
		SELECT DISTINCT nodecontrol.*
		FROM process_flow_control_node AS nodecontrol
		INNER JOIN process_flow_nodes AS node
			ON node.id = nodecontrol.node_id
		INNER JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		WHERE flow.id = $1
	`, flowId)
	return relationships, err
}

func FindRiskControlRelationships(orgId int32, role *core.Role) ([]*core.RiskControlRelationship, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	relationships := make([]*core.RiskControlRelationship, 0)
	err := dbConn.Select(&relationships, `
		SELECT DISTINCT riskcontrol.*
		FROM process_flow_risk_control AS riskcontrol
		INNER JOIN process_flow_risks AS risk
			ON risk.id = riskcontrol.risk_id
		INNER JOIN process_flow_controls AS control
			ON control.id = riskcontrol.control_id
		WHERE risk.org_id = $1
			AND control.org_id = $1
	`, orgId)
	return relationships, err
}

func FindControlsRelatedToRisk(riskId int64, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	controls := make([]*core.Control, 0)
	err := dbConn.Select(&controls, `
		SELECT DISTINCT control.*
		FROM process_flow_risk_control AS riskcontrol
		INNER JOIN process_flow_controls AS control
			ON control.id = riskcontrol.control_id
		WHERE riskcontrol.risk_id = $1
	`, riskId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, c := range controls {
		err = LogAuditSelectWithTx(c.OrgId, core.ResourceControl, strconv.FormatInt(c.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return controls, tx.Commit()
}

func FindRisksRelatedToControl(controlId int64, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	risks := make([]*core.Risk, 0)
	err := dbConn.Select(&risks, `
		SELECT DISTINCT risk.id, risk.name, risk.description
		FROM process_flow_risk_control AS riskcontrol
		INNER JOIN process_flow_risks AS risk
			ON risk.id = riskcontrol.risk_id
		WHERE riskcontrol.control_id = $1
	`, controlId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, r := range risks {
		err = LogAuditSelectWithTx(r.OrgId, core.ResourceRisk, strconv.FormatInt(r.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return risks, tx.Commit()
}
