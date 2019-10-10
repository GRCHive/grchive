package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

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
	err := dbConn.Select(&nodes, `
		SELECT * FROM process_flow_nodes WHERE process_flow_id = $1
	`, flowId)
	return nodes, err
}
