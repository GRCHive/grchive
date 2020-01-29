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

func GetAllVersionsFileStorage(fileId int64, orgId int32, role *core.Role) ([]*core.FileStorageData, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retData := make([]*core.FileStorageData, 0)
	err := dbConn.Select(&retData, `
		SELECT storage.*
		FROM file_storage AS storage
		INNER JOIN file_version_history AS fvh
			ON storage.id = fvh.file_storage_id
		WHERE fvh.file_id = $1
			AND fvh.org_id = $2
	`, fileId, orgId)
	return retData, err
}

func GetFileVersionStorageData(fileId int64, orgId int32, version int32, role *core.Role) (*core.FileStorageData, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := core.FileStorageData{}
	err := dbConn.Get(&data, `
		SELECT storage.*
		FROM file_storage AS storage
		INNER JOIN file_version_history AS fvh
			ON storage.id = fvh.file_storage_id
		WHERE fvh.file_id = $1
			AND fvh.org_id = $2
			AND fvh.version_number = $3
	`, fileId, orgId, version)
	return &data, err
}

func GetPreviewFileVersionStorageData(fileId int64, orgId int32, version int32, role *core.Role) (*core.FileStorageData, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT storage.*
		FROM file_storage AS storage
		INNER JOIN file_previews AS fp
			ON fp.preview_storage_id = storage.id
		INNER JOIN file_version_history AS fvh
			ON fp.file_id = fvh.file_id
		WHERE fvh.file_id = $1
			AND fvh.org_id = $2
			AND fvh.version_number = $3
			AND fp.preview_storage_id IS NOT NULL
	`, fileId, orgId, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	data := core.FileStorageData{}
	err = rows.StructScan(&data)
	if err != nil {
		return nil, err
	}
	return &data, err
}
