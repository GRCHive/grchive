package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func CreateFileStorageWithTx(storage *core.FileStorageData, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
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
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&storage.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFileStorageStorageIdWithTx(id int64, orgId int32, storageId string, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE file_storage
		SET storage_id = $1
		WHERE id = $2 AND org_id = $3	
	`, storageId, id, orgId)
	return err
}

func CreateControlDocumentationFileWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
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
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&file.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateControlDocumentationWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
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
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	err = UpdateControlDocumentationWithTx(file, tx, role)
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	for _, id := range fileIds {
		_, err = tx.Exec(`
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

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, r := range retArr {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdDocMetadata, strconv.FormatInt(r.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return retArr, tx.Commit()
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

	if err != nil {
		return nil, err
	}

	return &retFile, LogAuditSelect(orgId, core.ResourceIdDocMetadata, strconv.FormatInt(retFile.Id, 10), role)
}

func GetControlDocumentationStorage(fileId int64, orgId int32, role *core.Role) (*core.FileStorageData, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT storage.*
		FROM file_storage AS storage
		WHERE storage.metadata_id = $1 AND org_id = $2
	`, fileId, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	retFile := core.FileStorageData{}
	err = rows.StructScan(&retFile)
	if err != nil {
		return nil, err
	}

	return &retFile, LogAuditSelect(orgId, core.ResourceIdFileStorage, strconv.FormatInt(retFile.Id, 10), role)
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, r := range retArr {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdDocMetadata, strconv.FormatInt(r.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return retArr, tx.Commit()
}

func LinkFileWithPreviewWithTx(
	file core.ControlDocumentationFile,
	storage core.FileStorageData,
	preview core.FileStorageData,
	role *core.Role,
	tx *sqlx.Tx,
) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO file_previews (file_id, original_storage_id, preview_storage_id, org_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (original_storage_id)
			DO UPDATE
				SET preview_storage_id = EXCLUDED.preview_storage_id
	`, file.Id, storage.Id, preview.Id, file.OrgId)
	return err
}

func MarkPreviewUnavailable(file core.ControlDocumentationFile, storage core.FileStorageData, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO file_previews (file_id, original_storage_id, preview_storage_id, org_id)
		VALUES ($1, $2, NULL, $3)
		ON CONFLICT (original_storage_id)
			DO UPDATE
				SET preview_storage_id = EXCLUDED.preview_storage_id
	`, file.Id, storage.Id, file.OrgId)
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
		INNER JOIN file_comment_threads AS fc
			ON fc.thread_id = t.id
		WHERE fc.file_id = $1
			AND fc.org_id = $2
	`, fileId, orgId)
}

func GetDocumentCommentThreadId(fileId int64, orgId int32, role *core.Role) (int64, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return -1, core.ErrorUnauthorized
	}

	threadId := int64(-1)
	err := dbConn.Get(&threadId, `
		SELECT thread_id
		FROM file_comment_threads
		WHERE file_id = $1 AND org_id = $2
	`, fileId, orgId)
	return threadId, err
}

func InsertDocumentComment(fileId int64, orgId int32, comment *core.Comment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	threadId := int64(0)
	err := dbConn.Get(&threadId, `
		SELECT thread_id
		FROM file_comment_threads
		WHERE file_id = $1 AND org_id = $2
	`, fileId, orgId)
	if err != nil {
		return err
	}

	tx := dbConn.MustBegin()
	err = insertCommentWithTx(comment, threadId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetFileStorageAuxData(storageId int64, orgId int32, role *core.Role) (*core.FileStorageAuxData, error) {
	aux := core.FileStorageAuxData{}
	err := dbConn.Get(&aux, `
		SELECT 
			fp.file_id IS NOT NULL AS "is_preview",
			COALESCE(fvhp.version_number, fvhs.version_number, 0) AS "version_number"
		FROM file_storage AS s
		LEFT JOIN file_previews AS fp
			ON fp.preview_storage_id = s.id
		LEFT JOIN file_previews AS fs
			ON fs.original_storage_id = s.id
		LEFT JOIN file_version_history AS fvhp
			ON fvhp.file_storage_id = fp.original_storage_id
		LEFT JOIN file_version_history AS fvhs
			ON fvhs.file_storage_id = fs.original_storage_id
		WHERE s.id = $1 AND s.org_id = $2
	`, storageId, orgId)
	return &aux, err
}
