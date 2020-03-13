package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strconv"
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

func FindAllControlsForOrganization(org *core.Organization, filter core.ControlFilterData, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	controls := make([]*core.Control, 0)

	err := dbConn.Select(&controls, fmt.Sprintf(`
		SELECT control.*
		FROM process_flow_controls as control
		LEFT JOIN process_flow_risk_control AS riskcontrol
			ON riskcontrol.control_id = control.id
		WHERE control.org_id = $1
		GROUP BY control.id
		HAVING
			%s
		ORDER BY name ASC
	`,
		buildNumericFilter("COUNT(riskcontrol.risk_id)", filter.NumRisks),
	), org.Id)

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, c := range controls {
		err = LogAuditSelectWithTx(c.OrgId, core.ResourceControl, strconv.FormatInt(c.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return controls, tx.Commit()
}

func EditControl(control *core.Control, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE process_flow_controls
		SET name = :name, 
			description = :description,
			identifier = :identifier,
			control_type = :control_type,
			freq_type = :freq_type,
			freq_interval = :freq_interval,
			freq_other = :freq_other,
			owner_id = :owner_id,
			is_manual = :is_manual
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	if control.OwnerId.Valid {
		rows, err = tx.NamedQuery(`
			INSERT INTO process_flow_controls (name, identifier, description, control_type, org_id, freq_type, freq_interval, owner_id, is_manual, freq_other)
			VALUES (:name, :identifier, :description, :control_type, :org_id, :freq_type, :freq_interval, :owner_id, :is_manual, :freq_other)
			RETURNING id
		`, control)
	} else {
		rows, err = tx.NamedQuery(`
			INSERT INTO process_flow_controls (name, identifier, description, control_type, org_id, freq_type, freq_interval, is_manual, freq_other)
			VALUES (:name, :identifier, :description, :control_type, :org_id, :freq_type, :freq_interval, :is_manual, :freq_other)
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
		rows.Close()
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

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
			// Delete folders first because otherwise
			// the CASCADE from deleting process_flow_controls
			// will destroy the connection we have betwen the folder
			// and control.
			_, err := tx.Exec(`
				DELETE FROM file_folders AS folder
				USING control_folder_link AS link
				WHERE folder.id = link.folder_id
					AND link.control_id = $1
					AND link.org_id = $2
			`, id, orgId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = tx.Exec(`
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
	if err != nil {
		return nil, err
	}
	return &control, LogAuditSelect(control.OrgId, core.ResourceControl, strconv.FormatInt(control.Id, 10), role)
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
