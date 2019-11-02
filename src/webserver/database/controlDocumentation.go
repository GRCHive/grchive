package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"math"
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

func EditControlDocumentationCategory(cat *core.ControlDocumentationCategory) error {
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE process_flow_control_documentation_categories
		SET name = :name, 
			description = :description
		WHERE id = :id
	`, cat)
	if err != nil {
		tx.Rollback()
		return err
	}
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

func DeleteControlDocumentationCategory(catId int64) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM process_flow_control_documentation_categories
		WHERE id = $1
	`, catId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func CreateControlDocumentationFileWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_control_documentation_file (storage_name, relevant_time, upload_time, category_id)
		VALUES (:storage_name, :relevant_time, :upload_time, :category_id)
		RETURNING id
	`, file)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&file.Id)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func UpdateControlDocumentation(file *core.ControlDocumentationFile, tx *sqlx.Tx) error {
	_, err := tx.NamedExec(`
		UPDATE process_flow_control_documentation_file
		SET bucket_id = :bucket_id,
			storage_id = :storage_id
		WHERE id = :id
	`, file)
	return err
}

func GetControlDocumentation(catId int64, pageSize int, pageOffset int) ([]*core.ControlDocumentationFile, error) {
	retArr := make([]*core.ControlDocumentationFile, 0)

	err := dbConn.Select(&retArr, `
		SELECT *
		FROM process_flow_control_documentation_file
		WHERE category_id = $1
			AND bucket_id IS NOT NULL
			AND storage_id IS NOT NULL
		ORDER BY relevant_time DESC
		LIMIT $2
		OFFSET $3
	`, catId, pageSize, pageOffset)

	return retArr, err
}

func GetTotalControlDocumentationPages(catId int64, pageSize int) (int, error) {
	count := 0

	err := dbConn.Get(&count, `
		SELECT COUNT(*)
		FROM process_flow_control_documentation_file
		WHERE category_id = $1
			AND bucket_id IS NOT NULL
			AND storage_id IS NOT NULL
	`, catId)

	return int(math.Ceil(float64(count) / float64(pageSize))), err
}
