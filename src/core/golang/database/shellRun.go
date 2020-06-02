package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"time"
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

func GetShellRunServers(runId int64) ([]*core.Server, error) {
	servers := make([]*core.Server, 0)
	err := dbConn.Select(&servers, `
		SELECT s.*
		FROM infrastructure_servers AS s
		INNER JOIN shell_script_run_servers AS srs
			ON srs.server_id = s.id
		WHERE srs.run_id = $1
	`, runId)
	return servers, err
}

func GetShellRun(runId int64) (*core.ShellScriptRun, error) {
	run := core.ShellScriptRun{}
	err := dbConn.Get(&run, `
		SELECT *
		FROM shell_script_runs
		WHERE id = $1
	`, runId)
	return &run, err
}

func GetServerShellRuns(runId int64) ([]*core.ShellScriptRunPerServer, error) {
	run := []*core.ShellScriptRunPerServer{}
	err := dbConn.Select(&run, `
		SELECT *
		FROM shell_script_run_servers
		WHERE run_id = $1
	`, runId)
	return run, err
}

func GetShellRunForServer(runId int64, serverId int64) (*core.ShellScriptRunPerServer, error) {
	run := core.ShellScriptRunPerServer{}
	err := dbConn.Get(&run, `
		SELECT *
		FROM shell_script_run_servers
		WHERE run_id = $1 AND server_id = $2
	`, runId, serverId)
	return &run, err
}

func MarkShellScriptRunStartWithTx(tx *sqlx.Tx, runId int64, tm time.Time) error {
	_, err := tx.Exec(`
		UPDATE shell_script_runs
		SET run_time = $2
		WHERE id = $1
	`, runId, tm)
	return err
}

func MarkShellScriptRunEndWithTx(tx *sqlx.Tx, runId int64, tm time.Time) error {
	_, err := tx.Exec(`
		UPDATE shell_script_runs
		SET end_time = $2
		WHERE id = $1
	`, runId, tm)
	return err
}

func MarkShellScriptRunForServerStartWithTx(tx *sqlx.Tx, runId int64, serverId int64, tm time.Time) error {
	_, err := tx.Exec(`
		UPDATE shell_script_run_servers
		SET run_time = $3
		WHERE run_id = $1
			AND server_id = $2
	`, runId, serverId, tm)
	return err
}

func MarkShellScriptRunForServerFinishWithTx(tx *sqlx.Tx, runId int64, serverId int64, tm time.Time, success bool, log string) error {
	_, err := tx.Exec(`
		UPDATE shell_script_run_servers
		SET end_time = $3,
			success = $4,
			encrypted_log = $5
		WHERE run_id = $1
			AND server_id = $2
	`, runId, serverId, tm, success, log)
	return err
}

func AllShellRunsForServer(serverId int64) ([]*core.ShellScriptRun, error) {
	runs := make([]*core.ShellScriptRun, 0)
	err := dbConn.Select(&runs, `
		SELECT DISTINCT(ssr.*)
		FROM shell_script_runs AS ssr
		INNER JOIN shell_script_run_servers AS ssrs
			ON ssrs.run_id = ssr.id
		WHERE ssrs.serer_id = $1
		ORDER BY ssr.id DESC
	`, serverId)
	return runs, err
}

func AllShellRunsForShellScript(scriptId int64) ([]*core.ShellScriptRun, error) {
	runs := make([]*core.ShellScriptRun, 0)
	err := dbConn.Select(&runs, `
		SELECT ssr.*
		FROM shell_script_runs AS ssr
		INNER JOIN shell_script_versions AS ssv
			ON ssv.id = ssr.script_version_id
		WHERE ssv.shell_id = $1
		ORDER BY ssr.id DESC
	`, scriptId)
	return runs, err
}
