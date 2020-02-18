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
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&comment.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCommentWithTx(comment *core.Comment, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		UPDATE comments
		SET content = :content
		WHERE id = :id
			AND user_id = :user_id
		RETURNING *
	`, comment)

	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(comment)
	if err != nil {
		return err
	}
	return nil
}

func UpdateComment(comment *core.Comment) error {
	tx := dbConn.MustBegin()
	err := UpdateCommentWithTx(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteComment(commentId int64, userId int64) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM comments
		WHERE id = $1
			AND user_id = $2
	`, commentId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
