package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func NewControlDocumentationCategory(cat *core.ControlDocumentationCategory) error {
	tx := dbConn.MustBegin()

	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_control_documentation_categories (name, description, control_id)
		VALUES (:name, :description, :control_id)
		RETURNING id
	`, cat)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&cat.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func FindControlDocumentCategoriesForControl(controlId int64) ([]*core.ControlDocumentationCategory, error) {
	retArr := make([]*core.ControlDocumentationCategory, 0)

	err := dbConn.Select(&retArr, `
		SELECT cat.*
		FROM process_flow_control_documentation_categories AS cat
		WHERE cat.control_id = $1
	`, controlId)
	return retArr, err
}
