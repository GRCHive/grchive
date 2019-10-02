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
		return err
	}
	err = tx.Commit()
	return err
}
