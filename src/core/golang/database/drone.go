package database

func StartDroneCIJob(commitHash string) error {
	tx := CreateTx()

	_, err := tx.Exec(`
		INSERT INTO managed_code_drone_ci (
			code_id,
			org_id,
			commit_hash,
			time_start
		)
		SELECT id, org_id, git_hash, NOW()
		FROM managed_code
		WHERE git_hash = $1
		ON CONFLICT (commit_hash) DO UPDATE SET
			time_start = EXCLUDED.time_start,
			time_end = NULL,
			success = false,
			logs = '',
			jar = ''
	`, commitHash)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func FinishDroneCIJob(commitHash string, success bool, logs string, jar string) error {
	tx := CreateTx()

	_, err := tx.Exec(`
		UPDATE managed_code_drone_ci
		SET time_end = NOW(),
			success = $2,
			logs = $3,
			jar = $4
		WHERE commit_hash = $1
	`, commitHash, success, logs, jar)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
