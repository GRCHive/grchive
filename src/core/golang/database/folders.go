package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewFolderWithTx(folder *core.FileFolder, role *core.Role, tx *sqlx.Tx) error {
	// Better permission?
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO file_folders (org_id, name)
		VALUES (:org_id, :name)
		RETURNING id
	`, folder)

	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	err = rows.Scan(&folder.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFolder(folder *core.FileFolder, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE file_folders
		SET name = :name
		WHERE id = :id
			AND org_id = :org_id
	`, folder)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DeleteFolder(folderId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM file_folders
		WHERE id = $1
			AND org_id = $2
	`, folderId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
