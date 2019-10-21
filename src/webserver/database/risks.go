package database

import (
	"errors"
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"strings"
)

func InsertNewRisk(risk *core.RiskForNode) error {
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
