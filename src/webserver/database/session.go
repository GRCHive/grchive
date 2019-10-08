package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func StoreUserSession(session *core.UserSession) error {
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		INSERT INTO user_sessions
		VALUES (:session_id
			, :email
			, :last_active_time
			, :expiration_time
			, :user_agent
			, :ip_address
			, :access_token
			, :id_token
			, :refresh_token)
	`, session)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func FindUserSession(sessionId string) (*core.UserSession, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM user_sessions WHERE session_id = $1
	`, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var session *core.UserSession = new(core.UserSession)
	if !rows.Next() {
		return nil, rows.Err()
	}
	err = rows.StructScan(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func DeleteUserSession(sessionId string) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM user_sessions
		WHERE session_id = $1
	`, sessionId)
	if err != nil {
		return tx.Rollback()
	}
	err = tx.Commit()
	return err
}

func UpdateUserSession(session *core.UserSession) error {
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE user_sessions
		SET last_active_time = :last_active_time,
			expiration_time = :expiration_time,
			user_agent = :user_agent,
			ip_address = :ip_address,
			access_token = :access_token,
			id_token = :id_token,
			refresh_token = :refresh_token
		WHERE session_id = :session_id
	`, session)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
