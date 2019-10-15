package database

import (
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func GetAllProcessFlowIOTypes() ([]*core.ProcessFlowIOType, error) {
	result := []*core.ProcessFlowIOType{}

	err := dbConn.Select(&result, `
		SELECT * FROM process_flow_input_output_type ORDER BY name ASC`)
	return result, err
}

func CreateNewProcessFlowIO(io *core.ProcessFlowInputOutput, isInput bool) (*core.ProcessFlowInputOutput, error) {
	var err error
	var dbName string = ""
	if isInput {
		dbName = "process_flow_node_inputs"
	} else {
		dbName = "process_flow_node_outputs"
	}

	tx := dbConn.MustBegin()
	rows, err := tx.Queryx(fmt.Sprintf(`
		INSERT INTO %s (name, parent_node_id, io_type_id)
		VALUES ($1, $2, $3)
		RETURNING *
	`, dbName), io.Name, io.ParentNodeId, io.TypeId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	outIo := core.ProcessFlowInputOutput{}
	rows.Next()
	err = rows.StructScan(&outIo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	err = tx.Commit()
	return &outIo, err
}

func DeleteProcessFlowIO(ioId int64, isInput bool) error {
	var dbName string = ""
	if isInput {
		dbName = "process_flow_node_inputs"
	} else {
		dbName = "process_flow_node_outputs"
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, dbName), ioId)
	if err != nil {
		return tx.Rollback()
	}
	err = tx.Commit()
	return err
}
