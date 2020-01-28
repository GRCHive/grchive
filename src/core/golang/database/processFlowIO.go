package database

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
)

func getProcessFlowIODbName(isInput bool) string {
	if isInput {
		return "process_flow_node_inputs"
	} else {
		return "process_flow_node_outputs"
	}
}

func GetAllProcessFlowIOTypes(role *core.Role) ([]*core.ProcessFlowIOType, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	result := []*core.ProcessFlowIOType{}

	err := dbConn.Select(&result, `
		SELECT * FROM process_flow_input_output_type ORDER BY name ASC`)
	return result, err
}

func CreateNewProcessFlowIO(io *core.ProcessFlowInputOutput, isInput bool, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return nil, core.ErrorUnauthorized
	}
	var err error
	var dbName string = getProcessFlowIODbName(isInput)

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

func DeleteProcessFlowIO(ioId int64, isInput bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	var dbName string = getProcessFlowIODbName(isInput)
	tx := dbConn.MustBegin()
	_, err := tx.Exec(fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, dbName), ioId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func EditProcessFlowIO(io *core.ProcessFlowInputOutput, isInput bool, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return nil, core.ErrorUnauthorized
	}
	var dbName string = getProcessFlowIODbName(isInput)
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(fmt.Sprintf(`
		UPDATE %s
		SET name = :name, io_type_id = :io_type_id
		WHERE id = :id
		RETURNING *
	`, dbName), io)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	rows.Next()
	err = rows.StructScan(io)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	err = tx.Commit()
	return io, nil
}
