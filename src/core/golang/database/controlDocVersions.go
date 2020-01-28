package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AddFileVersionWithTx(file *core.ControlDocumentationFile, storage *core.FileStorageData, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO file_version_history (file_id, file_storage_id, org_id, version_number)
		SELECT $1, $2, $3, COALESCE(MAX(version_number), 0) + 1
		FROM file_version_history
		WHERE file_id = $1
	`, file.Id, storage.Id, file.OrgId)
	if err != nil {
		return err
	}

	return nil
}

func AllFileVersions(fileId int64, orgId int32, role *core.Role) ([]core.FileVersion, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	versions := make([]core.FileVersion, 0)
	err := dbConn.Select(&versions, `
		SELECT *
		FROM file_version_history
		WHERE file_id = $1
			AND org_id = $2
		ORDER BY version_number DESC
	`, fileId, orgId)
	return versions, err
}
