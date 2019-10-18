package database

import (
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"strconv"
)

type FlowNodeDisplaySettings map[string]interface{}

func FindDisplaySettingsForProcessFlow(flowId int64) (map[int64]FlowNodeDisplaySettings, error) {
	type QueryResult struct {
		NodeId   int64          `db:"parent_node_id"`
		Settings types.JSONText `db:"settings"`
	}

	// Read into an array and convert it into the map to make it easy to grab the
	// node id using data structures we already have.
	allResults := make([]QueryResult, 0)

	err := dbConn.Select(&allResults, `
		SELECT disp.parent_node_id, disp.settings
		FROM process_flow_node_display_settings AS disp
		LEFT JOIN process_flow_nodes AS node
			ON node.id = disp.parent_node_id
		LEFT JOIN process_flows AS flow
			ON flow.id = node.process_flow_id
		WHERE flow.id = $1
	`, flowId)

	if err != nil {
		return nil, err
	}

	retResult := make(map[int64]FlowNodeDisplaySettings)
	for _, result := range allResults {
		settings := FlowNodeDisplaySettings{}
		err = result.Settings.Unmarshal(&settings)
		if err != nil {
			return nil, err
		}
		retResult[result.NodeId] = settings
	}
	return retResult, nil
}

func FindDisplaySettingsForProcessFlowNode(nodeId int64) (FlowNodeDisplaySettings, error) {
	row := dbConn.QueryRowx(`
		SELECT settings 
		FROM process_flow_node_display_settings
		WHERE parent_node_id = $1
	`, nodeId)

	jsonData := types.JSONText{}
	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	retMap := map[string]interface{}{}
	err = jsonData.Unmarshal(retMap)
	if err != nil {
		return nil, err
	}

	return retMap, nil
}

func UpdateDisplaySettingsForProcessFlowNode(nodeId int64, settings map[string]interface{}) error {
	rawSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	var flowId int64 = 0

	tx := dbConn.MustBegin()
	rows, err := tx.Query(`
		WITH inserted AS (
			INSERT INTO process_flow_node_display_settings (parent_node_id, settings)
			VALUES ($1, $2)
			ON CONFLICT (parent_node_id)
				DO UPDATE 
					SET settings = EXCLUDED.settings
			RETURNING parent_node_id
		)
		SELECT flow.id
		FROM process_flow_nodes AS node
		INNER JOIN process_flows AS flow
			ON node.process_flow_id = flow.id
		WHERE node.id = $1
	`, nodeId, string(rawSettings))
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&flowId)
	if err != nil {
		return err
	}
	rows.Close()

	err = tx.Commit()
	if err == nil {
		core.SendMessage(core.UpdateDisplaySettingsForProcessFlowNode,
			core.MessageSubtype(strconv.FormatInt(flowId, 10)),
			struct {
				NodeId   int64
				Settings map[string]interface{}
			}{
				nodeId,
				settings,
			})
	}
	return err
}
