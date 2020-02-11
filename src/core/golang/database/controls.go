package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func GetControlTypes(role *core.Role) ([]*core.ControlType, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	retArr := make([]*core.ControlType, 0)

	err := dbConn.Select(&retArr, `
		SELECT * FROM process_flow_control_types ORDER BY name ASC`)

	return retArr, err
}

func FindAllControlsForOrganization(org *core.Organization, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	controls := make([]*core.Control, 0)

	err := dbConn.Select(&controls, `
		SELECT *
		FROM process_flow_controls as control
		WHERE control.org_id = $1
	`, org.Id)

	return controls, err
}

func EditControl(control *core.Control, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE process_flow_controls
		SET name = :name, 
			description = :description,
			control_type = :control_type,
			freq_type = :freq_type,
			freq_interval = :freq_interval,
			owner_id = :owner_id
		WHERE id = :id
	`, control)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func InsertNewControl(control *core.Control, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	var err error
	var rows *sqlx.Rows

	tx := dbConn.MustBegin()
	if control.OwnerId.Valid {
		rows, err = tx.NamedQuery(`
			INSERT INTO process_flow_controls (name, description, control_type, org_id, freq_type, freq_interval, owner_id)
			VALUES (:name, :description, :control_type, :org_id, :freq_type, :freq_interval, :owner_id)
			RETURNING id
		`, control)
	} else {
		rows, err = tx.NamedQuery(`
			INSERT INTO process_flow_controls (name, description, control_type, org_id, freq_type, freq_interval)
			VALUES (:name, :description, :control_type, :org_id, :freq_type, :freq_interval)
			RETURNING id
		`, control)
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&control.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func DeleteControls(nodeId int64, controlIds []int64, riskIds []int64, global bool, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessManage) ||
		!role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceRisks, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	// Always delete the control relationship between the node and the control
	// as well as the control and the specified risk.
	for idx, id := range controlIds {
		if nodeId != -1 {
			_, err := tx.Exec(`
				DELETE
				FROM process_flow_control_node AS cn
				USING process_flow_controls AS ctrl
				WHERE cn.control_id = ctrl.id
					AND cn.node_id = $1
					AND cn.control_id = $2
					AND ctrl.org_id = $3
			`, nodeId, id, orgId)

			if err != nil {
				tx.Rollback()
				return err
			}
		}

		if idx < len(riskIds) && riskIds[idx] != -1 {
			_, err := tx.Exec(`
				DELETE
				FROM process_flow_risk_control AS rc
				USING process_flow_controls AS ctrl
				WHERE rc.control_id = ctrl.id
					AND rc.risk_id = $1
					AND rc.control_id = $2
					AND ctrl.org_id = $3
			`, riskIds[idx], id, orgId)

			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// If global, delete the itself control (the relationships will be deleted by cascade).
	if global {
		for _, id := range controlIds {
			_, err := tx.Exec(`
				DELETE FROM process_flow_controls
				WHERE id = $1
					AND org_id = $2
			`, id, orgId)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func AddControlsToRisk(riskId int64, controlIds []int64, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	var err error
	tx := dbConn.MustBegin()

	for _, controlId := range controlIds {
		_, err = tx.Exec(`
			INSERT INTO process_flow_risk_control (risk_id, control_id)
			VALUES ($1, $2)
		`, riskId, controlId)
		// It's OK if we fail to add this if it's a duplicate.
		if err != nil && !IsDuplicateDBEntry(err) {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func AddControlsToNodeWithTx(nodeId int64, controlIds []int64, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	var err error
	for _, controlId := range controlIds {
		_, err = tx.Exec(`
			INSERT INTO process_flow_control_node (control_id, node_id)
			VALUES ($1, $2)
		`, controlId, nodeId)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddControlsToNode(nodeId int64, controlIds []int64, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := AddControlsToNodeWithTx(nodeId, controlIds, tx, role)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FindControl(controlId int64, role *core.Role) (*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	control := core.Control{}
	err := dbConn.Get(&control, `
		SELECT ctrl.*
		FROM process_flow_controls AS ctrl
		WHERE id = $1
	`, controlId)
	return &control, err
}

func getTableNameForControlDocCatIO(isInput bool) string {
	var tableName string
	if isInput {
		tableName = "controls_input_documentation"
	} else {
		tableName = "controls_output_documentation"
	}
	return tableName
}

func AddControlDocCatToControl(controlId int64, catId int64, orgId int32, isInput bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tableName := getTableNameForControlDocCatIO(isInput)

	tx := dbConn.MustBegin()
	_, err := tx.Exec(fmt.Sprintf(`
		INSERT INTO %s (category_id, org_id, control_id)
		VALUES ($1, $2, $3)
	`, tableName), catId, orgId, controlId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func RemoveControlDocCatFromControl(controlId int64, catId int64, orgId int32, isInput bool, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tableName := getTableNameForControlDocCatIO(isInput)

	tx := dbConn.MustBegin()
	_, err := tx.Exec(fmt.Sprintf(`
		DELETE FROM %s
		WHERE category_id = $1
			AND org_id = $2
			AND control_id = $3
	`, tableName), catId, orgId, controlId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
