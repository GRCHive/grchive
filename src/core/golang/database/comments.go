package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func getComments(condition string, args ...interface{}) ([]*core.Comment, error) {
	comments := make([]*core.Comment, 0)
	err := dbConn.Select(&comments, fmt.Sprintf(`
		SELECT comments.*
		FROM comments
		%s
		ORDER BY id DESC
	`, condition), args...)
	return comments, err
}

func insertCommentWithTx(comment *core.Comment, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO comments (
			user_id,
			post_time,
			content
		)
		VALUES (
			:user_id,
			:post_time,
			:content
		)
		RETURNING id
	`, comment)

	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&comment.Id)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
