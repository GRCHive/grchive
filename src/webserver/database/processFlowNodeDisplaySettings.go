package database

import (
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindDisplaySettingsForProcessFlowNode(nodeId int64) (map[string]interface{}, error) {
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

	tx := dbConn.MustBegin()
	_, err = tx.Exec(`
		INSERT INTO process_flow_node_display_settings (parent_node_id, settings)
		VALUES ($1, $2)
		ON CONFLICT (parent_node_id)
			DO UPDATE 
				SET settings = EXCLUDED.settings
	`, nodeId, string(rawSettings))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		core.SendMessage(core.UpdateDisplaySettingsForProcessFlowNode, struct {
			NodeId   int64
			Settings map[string]interface{}
		}{
			nodeId,
			settings,
		})
	}
	return err
}
