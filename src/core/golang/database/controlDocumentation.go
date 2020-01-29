package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func CreateFileStorageWithTx(storage *core.FileStorageData, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO file_storage (
			metadata_id,
			storage_name,
			org_id,
			bucket_id,
			storage_id,
			upload_time,
			upload_user_id
		)
		VALUES (
			:metadata_id,
			:storage_name,
			:org_id,
			:bucket_id,
			:storage_id,
			:upload_time,
			:upload_user_id
		)
		RETURNING id
	`, storage)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&storage.Id)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func CreateControlDocumentationFileWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO file_metadata (
			relevant_time,
			category_id,
			org_id,
			alt_name,
			description
		)
		VALUES (
			:relevant_time,
			:category_id,
			:org_id,
			:alt_name,
			:description
		)
		RETURNING id
	`, file)
	if err != nil {
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

func UpdateControlDocumentationWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.NamedExec(`
		UPDATE file_metadata
		SET alt_name = :alt_name,
			description = :description,
			relevant_time = :relevant_time,
			category_id = :category_id
		WHERE id = :id
			AND org_id = :org_id
	`, file)
	return err
}

func UpdateControlDocumentation(file *core.ControlDocumentationFile, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := UpdateControlDocumentationWithTx(file, tx, role)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteBatchControlDocumentation(fileIds []int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	for _, id := range fileIds {
		_, err := tx.Exec(`
			DELETE FROM file_metadata AS file
			WHERE file.id = $1
				AND file.org_id = $2
		`, id, orgId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetSocDocumentationForVendorProduct(productId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retArr := make([]*core.ControlDocumentationFile, 0)

	err := dbConn.Select(&retArr, `
		SELECT file.*
		FROM file_metadata as file
		INNER JOIN vendor_product_soc_reports AS soc
			ON soc.file_id = file.id
		WHERE soc.product_id = $1
			AND file.org_id = $2
		ORDER BY relevant_time DESC
	`, productId, orgId)

	return retArr, err

}

func GetControlDocumentation(fileId int64, orgId int32, role *core.Role) (*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retFile := core.ControlDocumentationFile{}

	err := dbConn.Get(&retFile, `
		SELECT file.*
		FROM file_metadata AS file
		INNER JOIN process_flow_control_documentation_categories AS cat
			ON file.category_id = cat.id
		WHERE file.id = $1
			AND file.org_id = $2
	`, fileId, orgId)

	return &retFile, err
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

func GetVersionedFileStorage(fileId int64, orgId int32, version int32, role *core.Role) (*core.FileStorageData, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retData := core.FileStorageData{}
	err := dbConn.Get(&retData, `
		SELECT storage.*
		FROM file_storage AS storage
		INNER JOIN file_version_history AS fvh
			ON storage.id = fvh.file_storage_id
		WHERE fvh.file_id = $1
			AND fvh.org_id = $2
			AND fvh.version_number = $3
	`, fileId, orgId, version)
	return &retData, err
}

func GetControlDocumentationPreview(fileId int64, orgId int32, role *core.Role) (*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT file.*
		FROM file_metadata AS file
		INNER JOIN file_previews AS pre
			ON file.id = pre.preview_file_id
		WHERE pre.original_file_id = $1
			AND pre.org_id = $2
	`, fileId, orgId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	retFile := core.ControlDocumentationFile{}
	err = rows.StructScan(&retFile)
	if err != nil {
		return nil, err
	}

	return &retFile, err
}

func GetControlDocumentationForCategory(catId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retArr := make([]*core.ControlDocumentationFile, 0)

	// LEFT JOIN on preview_file_id so that we don't
	// return files that are previews.
	err := dbConn.Select(&retArr, `
		SELECT file.*
		FROM file_metadata AS file
		WHERE file.category_id = $1
			AND file.org_id = $2
		ORDER BY file.relevant_time DESC
	`, catId, orgId)

	return retArr, err
}

func LinkFileWithPreviewWithTx(file core.ControlDocumentationFile, preview core.ControlDocumentationFile, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO file_previews (original_file_id, preview_file_id, category_id, org_id)
		VALUES ($1, $2, $3, $4)
	`, file.Id, preview.Id, file.CategoryId, file.OrgId)
	return err
}

func MarkPreviewUnavailable(file core.ControlDocumentationFile, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO file_previews (original_file_id, preview_file_id, category_id, org_id)
		VALUES ($1, NULL, $2, $3)
	`, file.Id, file.CategoryId, file.OrgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetDocumentComments(fileId int64, orgId int32, role *core.Role) ([]*core.Comment, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return getComments(`
		INNER JOIN file_comments AS fc
			ON fc.comment_id = comments.id
		WHERE fc.file_id = $1
			AND fc.org_id = $2
	`, fileId, orgId)
}

func InsertDocumentComment(fileId int64, catId int64, orgId int32, comment *core.Comment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	err := insertCommentWithTx(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO file_comments (
			file_id,
			cat_id,
			org_id,
			comment_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`, fileId, catId, orgId, comment.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
