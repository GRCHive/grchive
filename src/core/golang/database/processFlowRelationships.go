package database

import (
	"gitlab.com/grchive/grchive/core"
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

func FindNodesRelatedToRisk(riskId int64, role *core.Role) ([]*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	nodes := make([]*core.ProcessFlowNode, 0)
	err := dbConn.Select(&nodes, `
		SELECT DISTINCT node.*
		FROM process_flow_risk_node AS risknode
		INNER JOIN process_flow_nodes AS node
			ON node.id = risknode.node_id
		WHERE risknode.risk_id = $1
	`, riskId)
	return nodes, err
}

func FindNodesRelatedToControl(controlId int64, role *core.Role) ([]*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	nodes := make([]*core.ProcessFlowNode, 0)
	err := dbConn.Select(&nodes, `
		SELECT DISTINCT node.*
		FROM process_flow_control_node AS nodecontrol
		INNER JOIN process_flow_nodes AS node
			ON node.id = nodecontrol.node_id
		WHERE nodecontrol.control_id = $1
	`, controlId)
	return nodes, err
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
	return controls, err
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
	return risks, err
}
