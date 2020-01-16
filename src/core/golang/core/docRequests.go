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
