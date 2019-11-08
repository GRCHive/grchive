package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"time"
)

func StoreEmailVerification(veri core.EmailVerification) error {
	tx := dbConn.MustBegin()

	_, err := tx.NamedExec(`
		INSERT INTO email_verification (user_id, code, verification_sent)
		VALUES (:user_id, :code, :verification_sent)
		ON CONFLICT (user_id)
			DO UPDATE 
				SET code = EXCLUDED.code,
					verification_sent = EXCLUDED.verification_sent
	`, veri)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func FindUserVerification(code string, userId int64) (core.EmailVerification, error) {
	veri := core.EmailVerification{}

	err := dbConn.Get(&veri, `
		SELECT *
		FROM email_verification
		WHERE code = $1
			AND user_id = $2
	`, code, userId)

	return veri, err
}

func AcceptUserVerification(code string, userId int64) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		UPDATE email_verification
		SET verification_received = $1
		WHERE code = $2
			AND user_id = $3
	`, time.Now().UTC(), code, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func IsUserVerified(userId int64) (bool, error) {
	rows, err := dbConn.Queryx(`
		SELECT *
		FROM email_verification
		WHERE user_id = $1
			AND verification_received IS NOT NULL
	`, userId)

	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
