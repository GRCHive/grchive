package database

import (
	"gitlab.com/grchive/grchive/core"
)

func GetBuildLogs(hash string, orgId int32, role *core.Role) (string, error) {
	rows, err := dbConn.Queryx(`
		SELECT logs
		FROM managed_code_drone_ci
		WHERE commit_hash = $1 AND org_id = $2
	`, hash, orgId)

	if err != nil {
		return "", err
	}

	defer rows.Close()
	if !rows.Next() {
		return "", nil
	}

	logs := core.NullString{}
	err = rows.Scan(&logs)
	if err != nil {
		return "", err
	}

	if !logs.NullString.Valid {
		return "", nil
	}

	return logs.NullString.String, nil
}

func GetBuildLogsForRun(runId int64, role *core.Role) (string, error) {
	rows, err := dbConn.Queryx(`
		SELECT build_log
		FROM script_runs
		WHERE id = $1
	`, runId)

	if err != nil {
		return "", err
	}

	defer rows.Close()
	if !rows.Next() {
		return "", nil
	}

	logs := core.NullString{}
	err = rows.Scan(&logs)
	if err != nil {
		return "", err
	}

	if !logs.NullString.Valid {
		return "", nil
	}

	return logs.NullString.String, nil
}

func GetRunLogsForRun(runId int64, role *core.Role) (string, error) {
	rows, err := dbConn.Queryx(`
		SELECT run_log
		FROM script_runs
		WHERE id = $1
	`, runId)

	if err != nil {
		return "", err
	}

	defer rows.Close()
	if !rows.Next() {
		return "", nil
	}

	logs := core.NullString{}
	err = rows.Scan(&logs)
	if err != nil {
		return "", err
	}

	if !logs.NullString.Valid {
		return "", nil
	}

	return logs.NullString.String, nil
}
