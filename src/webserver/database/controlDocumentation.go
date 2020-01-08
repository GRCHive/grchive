package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"math"
)

func NewControlDocumentationCategory(cat *core.ControlDocumentationCategory, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_control_documentation_categories (name, description, org_id)
		VALUES (:name, :description, :org_id)
		RETURNING id
	`, cat)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&cat.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

func EditControlDocumentationCategory(cat *core.ControlDocumentationCategory, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE process_flow_control_documentation_categories
		SET name = :name, 
			description = :description
		WHERE id = :id
			AND org_id = :org_id
	`, cat)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteControlDocumentationCategory(catId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM process_flow_control_documentation_categories
		WHERE id = $1
			AND org_id = $2
	`, catId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func CreateControlDocumentationFileWithTx(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_control_documentation_file (
			storage_name,
			relevant_time,
			upload_time,
			category_id,
			org_id,
			upload_user_id,
			alt_name,
			description
		)
		VALUES (
			:storage_name,
			:relevant_time,
			:upload_time,
			:category_id,
			:org_id,
			:upload_user_id,
			:alt_name,
			:description
		)
		RETURNING id
	`, file)
	if err != nil {
		tx.Rollback()
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

func UpdateControlDocumentation(file *core.ControlDocumentationFile, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.NamedExec(`
		UPDATE process_flow_control_documentation_file
		SET bucket_id = :bucket_id,
			storage_id = :storage_id,
			alt_name = :alt_name,
			description = :description,
			upload_user_id = :upload_user_id,
			relevant_time = :relevant_time
		WHERE id = :id
			AND org_id = :org_id
			AND category_id = :category_id
	`, file)
	return err
}

func DeleteBatchControlDocumentation(fileIds []int64, catId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	for _, id := range fileIds {
		_, err := tx.Exec(`
			DELETE FROM process_flow_control_documentation_file AS file
			WHERE file.id = $1
				AND file.org_id = $2
				AND file.category_id = $3
		`, id, orgId, catId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetSocDocumentationForDeployment(deploymentId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retArr := make([]*core.ControlDocumentationFile, 0)

	err := dbConn.Select(&retArr, `
		SELECT file.*
		FROM process_flow_control_documentation_file as file
		INNER JOIN vendor_soc_reports AS soc
			ON soc.soc_report_file_id = file.id
				AND soc.soc_report_cat_id = file.category_id
		WHERE soc.deployment_id = $1
			AND file.org_id = $2
			AND bucket_id IS NOT NULL
			AND storage_id IS NOT NULL
		ORDER BY relevant_time DESC
	`, deploymentId, orgId)

	return retArr, err

}

func GetControlDocumentation(fileId int64, orgId int32, role *core.Role) (*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retFile := core.ControlDocumentationFile{}

	err := dbConn.Get(&retFile, `
		SELECT file.*
		FROM process_flow_control_documentation_file AS file
		INNER JOIN process_flow_control_documentation_categories AS cat
			ON file.category_id = cat.id
		WHERE file.id = $1
			AND file.org_id = $2
	`, fileId, orgId)

	return &retFile, err
}

func GetControlDocumentationForCategory(catId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationFile, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	retArr := make([]*core.ControlDocumentationFile, 0)

	err := dbConn.Select(&retArr, `
		SELECT *
		FROM process_flow_control_documentation_file
		WHERE category_id = $1
			AND org_id = $2
			AND bucket_id IS NOT NULL
			AND storage_id IS NOT NULL
		ORDER BY relevant_time DESC
	`, catId, orgId)

	return retArr, err
}

func GetTotalControlDocumentationPages(catId int64, orgId int32, pageSize int, role *core.Role) (int, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return 0, core.ErrorUnauthorized
	}

	count := 0

	err := dbConn.Get(&count, `
		SELECT COUNT(*)
		FROM process_flow_control_documentation_file
		WHERE category_id = $1
			AND org_id = $2
			AND bucket_id IS NOT NULL
			AND storage_id IS NOT NULL
	`, catId, orgId)

	return int(math.Ceil(float64(count) / float64(pageSize))), err
}

func GetAllDocumentationCategoriesForOrg(orgId int32, role *core.Role) ([]*core.ControlDocumentationCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	cats := make([]*core.ControlDocumentationCategory, 0)
	err := dbConn.Select(&cats, `
		SELECT *
		FROM process_flow_control_documentation_categories
		WHERE org_id = $1
	`, orgId)
	return cats, err
}

func GetDocumentationCategory(catId int64, orgId int32, role *core.Role) (*core.ControlDocumentationCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	cat := &core.ControlDocumentationCategory{}
	err := dbConn.Get(cat, `
		SELECT *
		FROM process_flow_control_documentation_categories
		WHERE id = $1
			AND org_id = $2
	`, catId, orgId)
	return cat, err
}

func getIoDocCatsForControl(table string, controlId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	cats := make([]*core.ControlDocumentationCategory, 0)
	err := dbConn.Select(&cats, fmt.Sprintf(`
		SELECT cat.*
		FROM %s AS cid
		INNER JOIN process_flow_control_documentation_categories AS cat
			ON cid.category_id = cat.id
		WHERE cid.control_id = $1
			AND cid.org_id = $2
	`, table), controlId, orgId)
	return cats, err
}

func GetInputDocumentationCategoriesForControl(controlId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationCategory, error) {
	return getIoDocCatsForControl("controls_input_documentation", controlId, orgId, role)
}

func GetOutputDocumentationCategoriesForControl(controlId int64, orgId int32, role *core.Role) ([]*core.ControlDocumentationCategory, error) {
	return getIoDocCatsForControl("controls_output_documentation", controlId, orgId, role)
}

func getControlsWithIoDocCat(table string, catId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	controls := make([]*core.Control, 0)
	err := dbConn.Select(&controls, fmt.Sprintf(`
		SELECT ctrl.*
		FROM %s AS cid
		INNER JOIN process_flow_controls AS ctrl
			ON cid.control_id = ctrl.id
		WHERE cid.category_id = $1
			AND cid.org_id = $2
	`, table), catId, orgId)
	return controls, err
}

func GetControlsWithInputDocumentationCategory(catId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	return getControlsWithIoDocCat("controls_input_documentation", catId, orgId, role)
}

func GetControlsWithOutputDocumentationCategory(catId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	return getControlsWithIoDocCat("controls_output_documentation", catId, orgId, role)
}
