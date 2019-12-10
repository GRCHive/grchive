package core

import (
	"time"
)

type DocumentRequest struct {
	Id              int64     `db:"id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	CatId           int64     `db:"cat_id"`
	OrgId           int32     `db:"org_id"`
	RequestedUserId int64     `db:"requested_user_id"`
	CompletionTime  NullTime  `db:"completion_time"`
	RequestTime     time.Time `db:"request_time"`
}

type DocumentRequestComment struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Text      string    `db:"text"`
	PostTime  time.Time `db:"post_time"`
	CatId     int64     `db:"cat_id"`
	OrgId     int32     `db:"org_id"`
	RequestId int64     `db:"request_id"`
}
