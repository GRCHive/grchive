package database

import (
	"gitlab.com/grchive/grchive/core"
)

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
	`, controlId, orgId)
	return folders, err
}
