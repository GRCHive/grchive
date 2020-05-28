package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func GetShellScriptsOfTypeForOrganization(shellType int32, orgId int32) ([]*core.ShellScript, error) {
	scripts := make([]*core.ShellScript, 0)
	err := dbConn.Select(&scripts, `
		SELECT *
		FROM shell_scripts
		WHERE type_id = $1
			AND org_id = $2
	`, shellType, orgId)
	return scripts, err
}

func AllShellScriptVersions(shellId int64, orgId int32) ([]*core.ShellScriptVersion, error) {
	versions := make([]*core.ShellScriptVersion, 0)
	err := dbConn.Select(&versions, `
		SELECT *
		FROM shell_script_versions
		WHERE shell_id = $1
			AND org_id = $2
	`, shellId, orgId)
	return versions, err
}

func UpdateShellScriptGCSStorageWithTx(tx *sqlx.Tx, shellId int64, bucketId string, storageId string) error {
	_, err := tx.Exec(`
		UPDATE shell_scripts
		SET bucket_id = $2,
			storage_id = $3
		WHERE id = $1
	`, shellId, bucketId, storageId)
	return err
}

func NewShellScriptWithTx(tx *sqlx.Tx, shell *core.ShellScript) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO shell_scripts (org_id, type_id, name, description)
		VALUES (:org_id, :type_id, :name, :description)
		RETURNING id
	`, shell)

	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	return rows.Scan(&shell.Id)
}

func CreateShellScriptVersionWithTx(tx *sqlx.Tx, shellId int64, orgId int32, uploadUser int64, generation int64) error {
	_, err := tx.Exec(`
		INSERT INTO shell_script_versions (shell_id, org_id, upload_time, upload_user_id, gcs_generation)
		VALUES ($1, $2, NOW(), $3, $4)
	`, shellId, orgId, uploadUser, generation)
	return err
}

func GetShellScriptFromId(shellId int64) (*core.ShellScript, error) {
	script := core.ShellScript{}
	err := dbConn.Get(&script, `
		SELECT *
		FROM shell_scripts
		WHERE id = $1
	`, shellId)
	return &script, err
}

func DeleteShellScriptFromIdWithTx(tx *sqlx.Tx, shellId int64) error {
	_, err := tx.Exec(`
		DELETE FROM shell_scripts
		WHERE id = $1
	`, shellId)
	return err
}
