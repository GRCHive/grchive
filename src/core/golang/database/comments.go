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
		INNER JOIN comment_threads AS t
			ON t.id = comments.thread_id 
		%s
		ORDER BY id DESC
	`, condition), args...)
	return comments, err
}

func insertCommentWithTx(comment *core.Comment, threadId int64, tx *sqlx.Tx) error {
	rows, err := tx.Queryx(`
		INSERT INTO comments (
			user_id,
			post_time,
			content,
			thread_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id
	`, comment.UserId, comment.PostTime, comment.Content, threadId)

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
		SET content = :content,
			post_time = :post_time
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

func FindUsersInCommentThread(threadId int64) ([]*core.User, error) {
	users := make([]*core.User, 0)
	err := dbConn.Select(&users, `
		SELECT DISTINCT u.*
		FROM users AS u
		INNER JOIN comments AS c
			ON c.user_id = u.id
		INNER JOIN comment_threads AS t
			ON t.id = c.thread_id
		WHERE t.id = $1
	`, threadId)
	return users, err
}
