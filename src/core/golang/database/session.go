package database

import (
	"gitlab.com/grchive/grchive/core"
)

func StoreUserSession(session *core.UserSession) error {
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		INSERT INTO user_sessions
		VALUES (:session_id
			, :last_active_time
			, :expiration_time
			, :user_agent
			, :ip_address
			, :access_token
			, :id_token
			, :refresh_token
			, :user_id)
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
	rows.Next()
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
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

// Use the sessionId passed in to find the session in the database
// so that we can use the session ID in the object to update the database.
func UpdateUserSession(session *core.UserSession, sessionId string) error {
	type UpdateType struct {
		Session  *core.UserSession `db:"sess"`
		SearchId string            `db:"id"`
	}

	updateData := UpdateType{
		Session:  session,
		SearchId: sessionId,
	}

	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE user_sessions
		SET session_id = :sess.session_id,
			last_active_time = :sess.last_active_time,
			expiration_time = :sess.expiration_time,
			user_agent = :sess.user_agent,
			ip_address = :sess.ip_address,
			access_token = :sess.access_token,
			id_token = :sess.id_token,
			refresh_token = :sess.refresh_token
		WHERE session_id = :id
	`, updateData)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
