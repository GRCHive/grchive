package database

import (
	"errors"
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"strings"
)

func DeleteRisks(nodeId int64, riskIds []int64, global bool) error {
	if len(riskIds) == 0 {
		return nil
	}

	riskIdQuery := make([]string, 0)
	riskIdParams := append(make([]interface{}, 0), nodeId)
	for idx, id := range riskIds {
		riskIdParams = append(riskIdParams, id)
		riskIdQuery = append(riskIdQuery, fmt.Sprintf("$%d", idx+2))
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(fmt.Sprintf(`
		DELETE FROM process_flow_risk_node
		WHERE node_id = $1
			AND risk_id IN (%s)
	`, strings.Join(riskIdQuery, ",")), riskIdParams...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if global {
		for _, id := range riskIds {
			_, err := tx.Exec(`
				DELETE FROM process_flow_risks
				WHERE id = $1
			`, id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	return err
}

func InsertNewRisk(risk *core.Risk) error {
	if len(risk.RelevantNodeIds) == 0 {
		return errors.New("No relevant node IDs")
	}

	var err error

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_risks (name, description)
		VALUES (:name, :description)
		RETURNING id
	`, risk)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&risk.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	valueConstruct := make([]string, 0)
	params := append(make([]interface{}, 0), risk.Id)
	for i, id := range risk.RelevantNodeIds {
		valueConstruct = append(valueConstruct, fmt.Sprintf("($1, $%d)", i+2))
		params = append(params, id)
	}

	riskNodeInsertStmt := fmt.Sprintf(`
		INSERT INTO process_flow_risk_node (risk_id, node_id)
		VALUES %s`,
		strings.Join(valueConstruct, ","))

	_, err = tx.Exec(riskNodeInsertStmt, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func FindAllRisksForProcessFlow(flowId int64) ([]*core.Risk, error) {
	risks := []*core.Risk{}
	rows, err := dbConn.Queryx(`
		SELECT 
			risk.*,
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(DISTINCT node.id), null)) AS node
		FROM process_flow_risks as risk
		INNER JOIN process_flow_risk_node AS risknode
			ON risknode.risk_id = risk.id
		INNER JOIN process_flow_nodes AS node
			ON risknode.node_id = node.id
		WHERE node.process_flow_id = $1
		GROUP BY risk.id
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
		newRisk := core.Risk{}
		newRisk.Id = dataMap["id"].(int64)
		newRisk.Name = dataMap["name"].(string)
		newRisk.Description = dataMap["description"].(string)
		newRisk.RelevantNodeIds, err = readInt64Array(dataMap["node"].([]uint8))
		if err != nil {
			return nil, err
		}

		risks = append(risks, &newRisk)
	}

	return risks, nil
}
