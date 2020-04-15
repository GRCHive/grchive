package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func AllManagedCodeForDataId(dataId int64, orgId int32, role *core.Role) ([]*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := make([]*core.ManagedCode, 0)
	err := dbConn.Select(&code, `
		SELECT code.*
		FROM managed_code AS code
		INNER JOIN code_to_client_data_link AS link
			ON link.code_id = code.id
		WHERE link.data_id = $1 AND link.org_id = $2
		ORDER BY code.id DESC
	`, dataId, orgId)
	return code, err
}

func AllManagedCodeForScriptId(scriptId int64, orgId int32, role *core.Role) ([]*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := make([]*core.ManagedCode, 0)
	err := dbConn.Select(&code, `
		SELECT code.*
		FROM managed_code AS code
		INNER JOIN code_to_client_scripts_link AS link
			ON link.code_id = code.id
		WHERE link.script_id = $1 AND link.org_id = $2
		ORDER BY code.id DESC
	`, scriptId, orgId)
	return code, err
}

func CheckValidCodeDataLink(codeId int64, dataId int64, orgId int32, role *core.Role) (bool, error) {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return false, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT *
		FROM code_to_client_data_link AS link
		WHERE link.code_id = $1 
			AND link.data_id = $2
			AND link.org_id = $3
	`, codeId, dataId, orgId)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func LinkCodeToData(codeId int64, dataId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := CreateTx()
	_, err := tx.Exec(`
		INSERT INTO code_to_client_data_link (code_id, data_id, org_id)
		VALUES ($1, $2, $3)
	`, codeId, dataId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func CheckValidCodeScriptLink(codeId int64, scriptId int64, orgId int32, role *core.Role) (bool, error) {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return false, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT *
		FROM code_to_client_scripts_link AS link
		WHERE link.code_id = $1 
			AND link.script_id = $2
			AND link.org_id = $3
	`, codeId, scriptId, orgId)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func LinkCodeToScriptWithTx(scriptId int64, dataId int64, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientScripts, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO code_to_client_scripts_link (code_id, script_id, org_id)
		VALUES ($1, $2, $3)
	`, scriptId, dataId, orgId)
	return err
}

func LinkCodeToScript(scriptId int64, dataId int64, orgId int32, role *core.Role) error {
	tx := CreateTx()
	err := LinkCodeToScriptWithTx(scriptId, dataId, orgId, role, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetLatestCodeForData(dataId int64, orgId int32, role *core.Role) (*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.ManagedCode{}
	err := dbConn.Get(&code, `
		SELECT code.*
		FROM managed_code AS code
		INNER JOIN code_to_client_data_link AS link
			ON link.code_id = code.id
		WHERE link.data_id = $1 AND link.org_id = $2
		ORDER BY code.id DESC
		LIMIT 1
	`, dataId, orgId)
	return &code, err
}

func GetLatestCodeForScript(scriptId int64, orgId int32, role *core.Role) (*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.ManagedCode{}
	err := dbConn.Get(&code, `
		SELECT code.*
		FROM managed_code AS code
		INNER JOIN code_to_client_scripts_link AS link
			ON link.code_id = code.id
		WHERE link.script_id = $1 AND link.org_id = $2
		ORDER BY code.id DESC
		LIMIT 1
	`, scriptId, orgId)
	return &code, err
}

func GetCode(codeId int64, orgId int32, role *core.Role) (*core.ManagedCode, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	code := core.ManagedCode{}
	err := dbConn.Get(&code, `
		SELECT *
		FROM managed_code
		WHERE id = $1 AND org_id = $2
	`, codeId, orgId)
	return &code, err
}

func InsertManagedCodeWithTx(code *core.ManagedCode, role *core.Role, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO managed_code (org_id, git_hash, action_time, git_path, gitea_file_sha)
		VALUES (:org_id, :git_hash, :action_time, :git_path, :gitea_file_sha)
		RETURNING id
	`, code)

	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&code.Id)
	if err != nil {
		return err
	}

	return nil
}

func InsertManagedCode(code *core.ManagedCode, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := CreateTx()
	err := InsertManagedCodeWithTx(code, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetCodeBuildStatus(commit string, orgId int32, role *core.Role) (*core.CodeBuildStatus, error) {
	if !role.Permissions.HasAccess(core.ResourceManagedCode, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT 
			CASE
				WHEN time_end IS NULL THEN true
				ELSE false
			END,
			success
		FROM managed_code_drone_ci
		WHERE commit_hash = $1 AND org_id = $2
	`, commit, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	status := core.CodeBuildStatus{
		Pending: true,
	}

	if rows.Next() {
		err = rows.Scan(&status.Pending, &status.Success)
		if err != nil {
			return nil, err
		}
	}

	return &status, nil
}

func CreateScriptRun(codeId int64, orgId int32, scriptId int64, role *core.Role) (*core.ScriptRun, error) {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	run := core.ScriptRun{}
	err = WrapTx(tx, func() error {
		rows, err := tx.Queryx(`
			INSERT INTO script_runs (link_id, start_time)
			SELECT link.id, NOW()
			FROM code_to_client_scripts_link AS link
			WHERE link.code_id = $1 AND link.org_id = $2 AND link.script_id = $3
			RETURNING *
		`, codeId, orgId, scriptId)

		if err != nil {
			return err
		}

		defer rows.Close()
		rows.Next()
		return rows.StructScan(&run)
	})
	return &run, err
}

func FinishBuildScriptRun(runId int64, success bool, logs string) error {
	tx := CreateTx()
	return WrapTx(tx, func() error {
		_, err := tx.Exec(`
			UPDATE script_runs
			SET build_success = $1,
				build_finish_time = NOW(),
				build_log = $2
			WHERE id = $3
		`, success, logs, runId)
		return err
	})
}

func GetScriptRun(runId int64) (*core.ScriptRun, error) {
	run := core.ScriptRun{}
	err := dbConn.Get(&run, `
		SELECT *
		FROM script_runs
		WHERE id = $1
	`, runId)
	return &run, err
}
