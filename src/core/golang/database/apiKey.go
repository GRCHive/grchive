package database

import (
	"gitlab.com/grchive/grchive/core"
)

func StoreApiKey(key *core.ApiKey) error {
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO api_keys ( hashed_api_key, expiration_date)
		VALUES (
			:hashed_api_key,
			:expiration_date,
		)
		RETURNING id
	`, key)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&key.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO api_key_to_users (api_key_id, user_id)
		VALUES (
			$1,
			$2
		)
	`, key.Id, key.UserId)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}

	rows.Close()
	err = tx.Commit()
	return err
}

func UpdateApiKey(key *core.ApiKey) error {
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE api_keys
		SET hashed_api_key = :hashed_api_key,
			expiration_date = :expiration_date
		WHERE id = :id
	`, key)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

// This function will return an error if an error occurs BUT
// will return a nil error (and a nil pointer) if no api key
// is found.
func FindApiKeyForUser(userId int64) (*core.ApiKey, error) {
	rows, err := dbConn.Queryx(`
		SELECT key.*, lnk.user_id
		FROM api_keys AS key
		INNER JOIN api_key_to_users AS lnk
			ON lnk.api_key_id = key.id
		WHERE lnk.user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	key := core.ApiKey{}
	err = rows.StructScan(&key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func FindApiKey(hashedRawKey string) (*core.ApiKey, error) {
	key := core.ApiKey{}

	err := dbConn.Get(&key, `
		SELECT key.*, lnk.user_id
		FROM api_keys AS key
		INNER JOIN api_key_to_users AS lnk
			ON lnk.api_key_id = key.id
		WHERE hashed_api_key = $1
	`, hashedRawKey)

	return &key, err
}

func DeleteApiKeyForUserId(userId int64) error {
	_, err := dbConn.Exec(`
		DELETE FROM api_keys AS key
		USING api_key_to_users AS lnk
		WHERE lnk.api_key_id = key.id AND lnk.user_id = $1
	`, userId)
	return err
}
