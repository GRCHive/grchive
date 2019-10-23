package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func GetControlTypes() ([]*core.ControlType, error) {
	retArr := make([]*core.ControlType, 0)

	err := dbConn.Select(&retArr, `
		SELECT * FROM process_flow_control_types ORDER BY name ASC`)

	return retArr, err
}

func FindAllControlsForOrganization(org *core.Organization) ([]*core.Control, error) {
	controls := make([]*core.Control, 0)

	err := dbConn.Select(&controls, `
		SELECT *
		FROM process_flow_controls as control
		WHERE control.org_id = $1
	`, org.Id)

	return controls, err
}

func InsertNewControl(control *core.Control, nodeId int64, riskId int64) error {
	var err error

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_controls (name, description, control_type, org_id, freq_type, freq_interval, owner_id)
		VALUES (:name, :description, :control_type.id, :org.id, :freq_type, :freq_interval, :owner.id)
		RETURNING id
	`, control)
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