package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func DeleteProcessFlowNodeFromId(nodeId int64, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
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

func GetAllProcessFlowNodeTypes(role *core.Role) ([]*core.ProcessFlowNodeType, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	result := []*core.ProcessFlowNodeType{}

	err := dbConn.Select(&result, `
		SELECT * FROM process_flow_node_types ORDER BY name ASC`)

	return result, err
}

func CreateNewProcessFlowNodeWithTypeIdWithTx(typeId int32, flowId int64, tx *sqlx.Tx, role *core.Role) (*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return nil, core.ErrorUnauthorized
	}
	var err error

	err = UpgradeTxToAudit(tx, role)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Queryx(`
		INSERT INTO process_flow_nodes (process_flow_id, node_type, name, description)
		VALUES ($1, $2, 'Temporary Name', '')
		RETURNING *
	`, flowId, typeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	node := core.ProcessFlowNode{}
	rows.Next()
	err = rows.StructScan(&node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func CreateNewProcessFlowNodeWithTypeId(typeId int32, flowId int64, role *core.Role) (*core.ProcessFlowNode, error) {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	node, err := CreateNewProcessFlowNodeWithTypeIdWithTx(typeId, flowId, tx, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return node, tx.Commit()
}

func findNodesHelper(role *core.Role, condition string, args ...interface{}) ([]*core.ProcessFlowNode, error) {
	nodes := []*core.ProcessFlowNode{}
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT 
			node.*,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT inp.*), null)) AS inputs,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT out.*), null)) AS outputs,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT risk.id), null)) AS risks,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT control.id), null)) AS controls,
			(ARRAY_AGG(DISTINCT flow.org_id))[1] AS "org_id"
		FROM process_flow_nodes AS node
		LEFT JOIN process_flow_node_inputs AS inp
			ON inp.parent_node_id = node.id
		LEFT JOIN process_flow_node_outputs AS out
			ON out.parent_node_id = node.id
		LEFT JOIN process_flow_risk_node AS risknode
			ON risknode.node_id = node.id
		LEFT JOIN process_flow_risks AS risk
			ON risknode.risk_id = risk.id
		LEFT JOIN process_flow_control_node AS controlnode
			ON controlnode.node_id = node.id
		LEFT JOIN process_flow_controls AS control
			ON controlnode.control_id = control.id
		LEFT JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		%s
		GROUP BY node.id
	`, condition), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dataMap := make(map[string]interface{})
		err = rows.MapScan(dataMap)
		if err != nil {
			tx.Rollback()
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
			tx.Rollback()
			return nil, err
		}
		newNode.Outputs, err = readProcessFlowInputOutputArray(dataMap["outputs"].([]uint8))
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// THESE ARE CONVERTING BYTE ARRAYS. THE uint8 IS CORRECT.
		newNode.RiskIds, err = readInt64Array(dataMap["risks"].([]uint8))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		newNode.ControlIds, err = readInt64Array(dataMap["controls"].([]uint8))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		nodes = append(nodes, &newNode)

		orgId := int32(dataMap["org_id"].(int64))
		err = LogAuditSelectWithTx(orgId, core.ResourceFlowNode, strconv.FormatInt(newNode.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return nodes, tx.Commit()
}

func FindNodeFromId(nodeId int64, role *core.Role) (*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	nodes, err := findNodesHelper(role, "WHERE node.id = $1", nodeId)
	if err != nil {
		return nil, err
	}

	if len(nodes) != 1 {
		return nil, errors.New("Unexpected number of nodes with the same node id.")
	}

	return nodes[0], nil
}

func FindAllNodesForProcessFlow(flowId int64, role *core.Role) ([]*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return findNodesHelper(role, "WHERE node.process_flow_id = $1", flowId)
}

func EditProcessFlowNodeWithTx(node *core.ProcessFlowNode, tx *sqlx.Tx, role *core.Role) (*core.ProcessFlowNode, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return nil, core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return nil, err
	}

	rows, err := tx.NamedQuery(`
		UPDATE process_flow_nodes
		SET name = :name, description = :description, node_type = :node_type
		WHERE id = :id
		RETURNING *
	`, node)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func EditProcessFlowNode(node *core.ProcessFlowNode, role *core.Role) (*core.ProcessFlowNode, error) {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	_, err = EditProcessFlowNodeWithTx(node, tx, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return node, tx.Commit()
}
