package database

import (
	"github.com/lib/pq"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

// bool: Returns whether or not the inserted data is a duplicate.
// error: Returns the raw database error.
func AddNewGettingStartedInterest(name string, email string) (bool, error) {
	var err error

	tx := dbConn.MustBegin()
	_, err = tx.Exec(`
		INSERT INTO get_started_interest (name, email)
		VALUES ($1, $2)
	`, name, email)
	if err != nil {
		core.Info(err.Error())
		tx.Rollback()
		return err.(*pq.Error).Code == "23505", err
	}

	err = tx.Commit()
	if err != nil {
		core.Info(err.Error())
		return false, err
	}
	return false, nil
}
