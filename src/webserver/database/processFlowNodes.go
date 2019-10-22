package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func DeleteProcessFlowNodeFromId(nodeId int64) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM process_flow_nodes
		WHERE id = $1
	`, nodeId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func GetAllProcessFlowNodeTypes() ([]*core.ProcessFlowNodeType, error) {
	result := []*core.ProcessFlowNodeType{}

	err := dbConn.Select(&result, `
		SELECT * FROM process_flow_node_types ORDER BY name ASC`)

	return result, err
}

func CreateNewProcessFlowNodeWithTypeId(typeId int32, flowId int64) (*core.ProcessFlowNode, error) {
	var err error

	tx := dbConn.MustBegin()
	rows, err := tx.Queryx(`
		INSERT INTO process_flow_nodes (process_flow_id, node_type, name, description)
		VALUES ($1, $2, 'Temporary Name', '')
		RETURNING *
	`, flowId, typeId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	node := core.ProcessFlowNode{}
	rows.Next()
	err = rows.StructScan(&node)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	err = tx.Commit()
	return &node, err
}

func FindAllNodesForProcessFlow(flowId int64) ([]*core.ProcessFlowNode, error) {
	nodes := []*core.ProcessFlowNode{}
	rows, err := dbConn.Queryx(`
		SELECT 
			node.*,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT inp.*), null)) AS inputs,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT out.*), null)) AS outputs,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT risk.id), null)) AS risks
		FROM process_flow_nodes AS node
		LEFT JOIN process_flow_node_inputs AS inp
			ON inp.parent_node_id = node.id
		LEFT JOIN process_flow_node_outputs AS out
			ON out.parent_node_id = node.id
		LEFT JOIN process_flow_risk_node AS risknode
			ON risknode.node_id = node.id
		LEFT JOIN process_flow_risks AS risk
			ON risknode.risk_id = risk.id
		WHERE node.process_flow_id = $1
		GROUP BY node.id
	`, flowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dataMap := make(map[string]interface{})
		err = rows.MapScan(dataMap)
		if err != nil {
			return nil, err
		}

		// Manual unmarshaling so we can support the JSON inputs/outputs.
		newNode := core.ProcessFlowNode{}
		newNode.Id = dataMap["id"].(int64)
		newNode.Name = dataMap["name"].(string)
		newNode.ProcessFlowId = dataMap["process_flow_id"].(int64)
		newNode.Description = dataMap["description"].(string)
		newNode.NodeTypeId = int32(dataMap["node_type"].(int64))
		newNode.Inputs, err = readProcessFlowInputOutputArray(dataMap["inputs"].([]uint8))
		if err != nil {
			return nil, err
		}
		newNode.Outputs, err = readProcessFlowInputOutputArray(dataMap["outputs"].([]uint8))
		if err != nil {
			return nil, err
		}
		newNode.RiskIds, err = readInt64Array(dataMap["risks"].([]uint8))
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &newNode)
	}

	return nodes, err
}

func EditProcessFlowNode(node *core.ProcessFlowNode) (*core.ProcessFlowNode, error) {
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		UPDATE process_flow_nodes
		SET name = :name, description = :description, node_type = :node_type
		WHERE id = :id
		RETURNING *
	`, node)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	rows.Next()
	err = rows.StructScan(node)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	err = tx.Commit()
	return node, nil
}
