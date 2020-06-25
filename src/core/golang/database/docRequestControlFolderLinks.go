package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AddDocRequestControlFolderLinkWithTx(requestId int64, controlId int64, folderId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDocRequests, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceControls, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	// Be a little more roundabount to ensure that the control-folder link actually exists.
	_, err := tx.Exec(`
		INSERT INTO request_folder_link (request_id, folder_id, org_id)
		SELECT $4, lnk.folder_id, lnk.org_id
		FROM control_folder_link AS lnk
		WHERE lnk.control_id = $1 AND lnk.folder_id = $2 AND lnk.org_id = $3
	`, controlId, folderId, orgId, requestId)
	return err
}

func FindFolderLinkedToDocRequestControl(requestId int64, controlId int64, orgId int32) (*core.FileFolder, error) {
	folder := core.FileFolder{}
	err := dbConn.Get(&folder, `
		SELECT fol.*
		FROM file_folders AS fol
		INNER JOIN request_folder_link AS lnk
			ON lnk.folder_id = fol.id
		INNER JOIN control_folder_link AS ctrl
			ON ctrl.folder_id = fol.id
		WHERE lnk.request_id = $1 
			AND ctrl.control_id = $2
			AND fol.org_id = $3
	`, requestId, controlId, orgId)
	return &folder, err
}
