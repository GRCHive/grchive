package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func NewControlDocumentationCategoryWithTx(cat *core.ControlDocumentationCategory, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceControlDocumentationMetadata, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
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
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = NewControlDocumentationCategoryWithTx(cat, role, tx)
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
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

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
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
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = DeleteControlDocumentationCategoryWithTx(catId, orgId, role, tx)
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

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, c := range cats {
		err = LogAuditSelectWithTx(orgId, core.ResourceDocCat, strconv.FormatInt(c.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return cats, tx.Commit()
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
	if err != nil {
		return nil, err
	}
	return cat, LogAuditSelect(orgId, core.ResourceDocCat, strconv.FormatInt(cat.Id, 10), role)
}
