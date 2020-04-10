package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewClientScriptWithTx(script *core.ClientScript, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO client_scripts (org_id, name, description)
		VALUES (:org_id, :name, :description)
		RETURNING id
	`, script)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&script.Id)
	if err != nil {
		return err
	}

	return nil
}

func NewClientScript(script *core.ClientScript, role *core.Role) error {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = NewClientScriptWithTx(script, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateClientScript(data *core.ClientScript, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE client_scripts
		SET name = :name,
			description = :description
		WHERE id = :id
			AND org_id = :org_id
	`, data)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func AllClientScriptsForOrganization(orgId int32, role *core.Role) ([]*core.ClientScript, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	scripts := make([]*core.ClientScript, 0)
	err := dbConn.Select(&scripts, `
		SELECT *
		FROM client_scripts
		WHERE org_id = $1
	`, orgId)
	return scripts, err
}

func GetClientScriptFromId(scriptId int64, orgId int32, role *core.Role) (*core.ClientScript, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	script := core.ClientScript{}
	err := dbConn.Get(&script, `
		SELECT *
		FROM client_scripts
		WHERE id = $1 AND org_id = $2
	`, scriptId, orgId)
	return &script, err
}

func DeleteClientScript(scriptId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM client_scripts
		WHERE id = $1 AND org_id = $2
	`, scriptId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
