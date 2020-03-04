package database

import (
	"gitlab.com/grchive/grchive/core"
)

func GetOrgGLCategories(orgId int32, role *core.Role) ([]*core.GeneralLedgerCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	cats := make([]*core.GeneralLedgerCategory, 0)
	err := dbConn.Select(&cats, `
		SELECT *
		FROM general_ledger_categories
		WHERE org_id = $1
	`, orgId)

	return cats, err
}

func GetOrgGLAccounts(orgId int32, role *core.Role) ([]*core.GeneralLedgerAccount, error) {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	accs := make([]*core.GeneralLedgerAccount, 0)
	err := dbConn.Select(&accs, `
		SELECT *
		FROM general_ledger_accounts
		WHERE org_id = $1
	`, orgId)

	return accs, err
}

func CreateNewGLCategory(cat *core.GeneralLedgerCategory, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO general_ledger_categories(org_id, parent_category_id, name, description)
		VALUES (:org_id, :parent_category_id, :name, :description)
		RETURNING id
	`, cat)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&cat.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func UpdateGLCategory(cat *core.GeneralLedgerCategory, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE general_ledger_categories
		SET parent_category_id = :parent_category_id,
			name = :name,
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

func CreateNewGLAccount(acc *core.GeneralLedgerAccount, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO general_ledger_accounts(org_id, parent_category_id, account_identifier, account_name, account_description, financially_relevant)
		VALUES (:org_id, :parent_category_id, :account_identifier, :account_name, :account_description, :financially_relevant)
		RETURNING id
	`, acc)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&acc.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func DeleteGLCategory(catId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM general_ledger_categories
		WHERE id = $1
			AND org_id = $2
	`, catId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetGLAccountFromDbId(accId int64, orgId int32, role *core.Role) (*core.GeneralLedgerAccount, error) {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	acc := core.GeneralLedgerAccount{}
	err := dbConn.Get(&acc, `
		SELECT *
		FROM general_ledger_accounts
		WHERE id = $1
			AND org_id = $2
	`, accId, orgId)

	return &acc, err
}

func FindGLAccountParentCategories(acc *core.GeneralLedgerAccount, role *core.Role) ([]*core.GeneralLedgerCategory, error) {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	parents := make([]*core.GeneralLedgerCategory, 0)
	err := dbConn.Select(&parents, `
		WITH RECURSIVE parents AS (
			SELECT cat.*
			FROM general_ledger_categories AS cat
			WHERE cat.id = $1
				AND cat.org_id = $2
			UNION
				SELECT cat.*
				FROM general_ledger_categories AS cat
				INNER JOIN parents
					ON parents.parent_category_id = cat.id
		)
		SELECT * FROM parents
	`, acc.ParentCategoryId, acc.OrgId)
	return parents, err
}

func UpdateGLAccount(acc *core.GeneralLedgerAccount, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE general_ledger_accounts
		SET 
			parent_category_id = :parent_category_id,
			account_identifier = :account_identifier,
			account_name = :account_name,
			account_description = :account_description,
			financially_relevant = :financially_relevant
		WHERE id = :id
			AND org_id = :org_id
	`, acc)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteGLAccount(accId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM general_ledger_accounts
		WHERE id = $1
			AND org_id = $2
	`, accId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}
