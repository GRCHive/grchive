package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindNodeRiskRelationships(flowId int64) ([]*core.NodeRiskRelationship, error) {
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

func FindNodeControlRelationships(flowId int64) ([]*core.NodeControlRelationship, error) {
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

func FindRiskControlRelationships(orgId int32) ([]*core.RiskControlRelationship, error) {
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
