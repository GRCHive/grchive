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

func LinkScriptToParameterWithTx(scriptId int64, codeId int64, orgId int32, name string, paramId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO client_script_code_parameters(link_id, name, param_type)
		SELECT link.id, $1, $2
		FROM code_to_client_scripts_link AS link
		WHERE link.script_id = $3 AND link.org_id = $4 AND link.code_id = $5
	`, name, paramId, scriptId, orgId, codeId)
	return nil
}

func LinkScriptToDataSourceWithTx(scriptId int64, codeId int64, orgId int32, dataId int64, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO client_script_code_data_sources(link_id, data_id, org_id)
		SELECT link.id, $1, $2
		FROM code_to_client_scripts_link AS link
		WHERE link.script_id = $3 AND link.org_id = $2 AND link.code_id = $4
	`, dataId, orgId, scriptId, codeId)
	return err
}

func GetLinkedDataSourceToScriptCode(scriptId int64, codeId int64, orgId int32, role *core.Role) ([]*core.FullClientDataWithLink, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceClientData, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return getClientDataHelper(role, `
		INNER JOIN client_script_code_data_sources AS src
			ON src.data_id = data.id
		INNER JOIN code_to_client_scripts_link AS lnk
			ON lnk.id = src.link_id
		WHERE lnk.script_id = $1 AND lnk.code_id = $2 AND lnk.org_id = $3
	`, scriptId, codeId, orgId)
}

func GetLinkedParametersToScriptCode(scriptId int64, codeId int64, orgId int32, role *core.Role) ([]*core.CodeParameter, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	params := make([]*core.CodeParameter, 0)
	err := dbConn.Select(&params, `
		SELECT param.*
		FROM client_script_code_parameters AS param
		INNER JOIN code_to_client_scripts_link AS lnk
			ON lnk.id = param.link_id
		WHERE lnk.script_id = $1 AND lnk.code_id = $2 AND lnk.org_id = $3
	`, scriptId, codeId, orgId)
	return params, err
}

func GetScriptForCode(codeId int64, orgId int32, role *core.Role) (*core.ClientScript, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT script.*
		FROM client_scripts AS script
		INNER JOIN code_to_client_scripts_link AS link
			ON link.script_id = script.id
		WHERE link.code_id = $1 AND link.org_id = $2
	`, codeId, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	script := core.ClientScript{}
	err = rows.StructScan(&script)
	return &script, err
}

func GetScriptFromScriptCodeLink(linkId int64, role *core.Role) (*core.ClientScript, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	script := core.ClientScript{}
	err := dbConn.Get(&script, `
		SELECT script.*
		FROM client_scripts AS script
		INNER JOIN code_to_client_scripts_link AS link
			ON link.script_id = script.id
		WHERE link.id = $1
	`, linkId)
	return &script, err
}

func GetCodeFromScriptCodeLink(linkId int64, role *core.Role) (*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.ManagedCode{}
	err := dbConn.Get(&code, `
		SELECT code.*
		FROM managed_code AS code
		INNER JOIN code_to_client_scripts_link AS link
			ON link.code_id = code.id
		WHERE link.id = $1
	`, linkId)
	return &code, err
}
