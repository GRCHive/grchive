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
