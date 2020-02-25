package database

import (
	"gitlab.com/grchive/grchive/core"
)

func FindControlsLinkedToDocCat(catId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	controls := make([]*core.Control, 0)
	err := dbConn.Select(&controls, `
		SELECT DISTINCT ctrl.*
		FROM process_flow_controls AS ctrl
		INNER JOIN control_folder_link AS cf
			ON cf.control_id = ctrl.id
		INNER JOIN file_folder_link AS ff
			ON ff.folder_id = cf.folder_id
		INNER JOIN file_metadata AS file
			ON file.id = ff.file_id
		WHERE file.category_id = $1
			AND file.org_id = $2
	`, catId, orgId)
	return controls, err
}