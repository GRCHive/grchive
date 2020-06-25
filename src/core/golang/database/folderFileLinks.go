package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func DeleteFileFromFolder(fileId int64, folderId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM file_folder_link
		WHERE file_id = $1
			AND folder_id = $2
			AND org_id = $3
	`, fileId, folderId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func AddFileToFolderWithTx(fileId int64, folderId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO file_folder_link (file_id, folder_id, org_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (org_id, folder_id, file_id) DO NOTHING
	`, fileId, folderId, orgId)
	return err
}

func FindFilesLinkedToFolder(folderId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	files := make([]*core.ControlDocumentationFile, 0)
	err := dbConn.Select(&files, `
		SELECT file.*
		FROM file_metadata AS file
		INNER JOIN file_folder_link AS link
			ON link.file_id = file.id
		WHERE link.folder_id = $1 AND link.org_id = $2
	`, folderId, orgId)
	return files, err
}
