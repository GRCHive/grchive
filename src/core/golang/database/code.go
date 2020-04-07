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
