package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewControlDocumentationCategoryWithTx(cat *core.ControlDocumentationCategory, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_control_documentation_categories (name, description, org_id)
		VALUES (:name, :description, :org_id)
		RETURNING id
	`, cat)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&cat.Id)
	if err != nil {
		return err
	}
	return nil
}

func NewControlDocumentationCategory(cat *core.ControlDocumentationCategory, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := NewControlDocumentationCategoryWithTx(cat, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

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

func DeleteControlDocumentationCategoryWithTx(catId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	_, err := tx.Exec(`
		DELETE FROM process_flow_control_documentation_categories
		WHERE id = $1
			AND org_id = $2
	`, catId, orgId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteControlDocumentationCategory(catId int64, orgId int32, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := DeleteControlDocumentationCategoryWithTx(catId, orgId, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
