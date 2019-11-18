package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
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

	tx := dbConn.MustBegin()
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
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func CreateNewGLAccount(acc *core.GeneralLedgerAccount, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
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
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}
