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
