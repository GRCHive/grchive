package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewShellRunWithTx(tx *sqlx.Tx, run *core.ShellScriptRun) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO shell_script_runs (script_version_id, run_user_id, create_time)
		VALUES (:script_version_id, :run_user_id, :create_time)
		RETURNING id
	`, run)

	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	return rows.Scan(&run.Id)
}

func NewShellServerRunWithTx(tx *sqlx.Tx, run *core.ShellScriptRunPerServer) error {
	_, err := tx.NamedExec(`
		INSERT INTO shell_script_run_servers (run_id, org_id, server_id)
		VALUES (:run_id, :org_id, :server_id)
	`, run)
	return err
}
