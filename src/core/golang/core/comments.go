package core

import (
	"time"
)

type Comment struct {
	Id       int64     `db:"id"`
	UserId   int64     `db:"user_id"`
	PostTime time.Time `db:"post_time"`
	Content  string    `db:"content"`
}
