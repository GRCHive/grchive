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
