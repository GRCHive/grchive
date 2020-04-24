package core

import (
	"time"
)

type GenericRequest struct {
	Id           int64     `db:"id"`
	OrgId        int32     `db:"org_id"`
	UploadTime   time.Time `db:"upload_time"`
	UploadUserId int64     `db:"upload_user_id"`
	Name         string    `db:"name"`
	Assignee     NullInt64 `db:"assignee"`
	DueDate      NullTime  `db:"due_date"`
	Description  string    `db:"description"`
}

type GenericApproval struct {
	Id              int64     `db:"id"`
	RequestId       int64     `db:"request_id"`
	ResponseTime    time.Time `db:"response_time"`
	ResponderUserId int64     `db:"responder_user_id"`
	Response        bool      `db:"response"`
	Reason          string    `db:"reason"`
}
