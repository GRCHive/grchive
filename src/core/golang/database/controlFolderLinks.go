package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AddControlFolderLinkWithTx(controlId int64, folderId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO control_folder_link (control_id, folder_id, org_id)
		VALUES ($1, $2, $3)
	`, controlId, folderId, orgId)
	return err
}

func FindFoldersLinkedToControl(controlId int64, orgId int32, role *core.Role) ([]*core.FileFolder, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	folders := make([]*core.FileFolder, 0)
	err := dbConn.Select(&folders, `
		SELECT folder.*
		FROM file_folders AS folder
		INNER JOIN control_folder_link AS link
			ON link.folder_id = folder.id
		WHERE link.control_id = $1 AND link.org_id = $2
		ORDER BY name
	`, controlId, orgId)
	return folders, err
}
