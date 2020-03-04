package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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

func CreateNewProcessFlowIOWithTx(io *core.ProcessFlowInputOutput, isInput bool, tx *sqlx.Tx, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return nil, core.ErrorUnauthorized
	}
	var err error
	var dbName string = getProcessFlowIODbName(isInput)

	err = UpgradeTxToAudit(tx, role)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Queryx(fmt.Sprintf(`
		INSERT INTO %s (name, parent_node_id, io_type_id, io_order)
		SELECT $1, $2, $3, COALESCE(MAX(io_order), 0) + 1
		FROM %s
		WHERE parent_node_id = $2 AND io_type_id = $3
		RETURNING *
	`, dbName, dbName), io.Name, io.ParentNodeId, io.TypeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outIo := core.ProcessFlowInputOutput{}
	rows.Next()
	err = rows.StructScan(&outIo)
	if err != nil {
		return nil, err
	}
	return &outIo, nil
}

func CreateNewProcessFlowIO(io *core.ProcessFlowInputOutput, isInput bool, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}
	outIo, err := CreateNewProcessFlowIOWithTx(io, isInput, tx, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return outIo, tx.Commit()
}

func DeleteProcessFlowIO(ioId int64, isInput bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	var dbName string = getProcessFlowIODbName(isInput)

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

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
		rows.Close()
		tx.Rollback()
		return nil, err
	}
	rows.Close()

	err = tx.Commit()
	return io, nil
}

func GetProcessFlowIOFromId(ioId int64, isInput bool, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbName := getProcessFlowIODbName(isInput)
	io := core.ProcessFlowInputOutput{}
	err := dbConn.Get(&io, fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE id = $1
	`, dbName), ioId)
	return &io, err
}

func GetSwapProcessFlowIO(io *core.ProcessFlowInputOutput, isInput bool, dir int32, role *core.Role) (*core.ProcessFlowInputOutput, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbName := getProcessFlowIODbName(isInput)

	var orderBy string
	var where string
	if dir < 0 {
		orderBy = "DESC"
		where = "<"
	} else {
		orderBy = "ASC"
		where = ">"
	}

	retIo := core.ProcessFlowInputOutput{}
	err := dbConn.Get(&retIo, fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE parent_node_id = $1
			AND io_type_id = $2
			AND io_order %s $3
		ORDER BY io_order %s
		LIMIT 1
	`, dbName, where, orderBy), io.ParentNodeId, io.TypeId, io.IoOrder)
	return &retIo, err
}

func SwapIOOrder(a *core.ProcessFlowInputOutput, b *core.ProcessFlowInputOutput, isInput bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	dbName := getProcessFlowIODbName(isInput)

	_, err = tx.Exec(fmt.Sprintf(`
		UPDATE %s
		SET io_order = CASE id
					   		WHEN $1 THEN (SELECT io_order FROM %s WHERE id = $2)
					   		WHEN $2 THEN (SELECT io_order FROM %s WHERE id = $1)
					   END
		WHERE id in ($1, $2)
	`, dbName, dbName, dbName), a.Id, b.Id)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
