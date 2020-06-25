package core

import (
	"time"
)

type DocumentRequest struct {
	Id              int64     `db:"id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	OrgId           int32     `db:"org_id"`
	RequestedUserId int64     `db:"requested_user_id"`
	AssigneeUserId  NullInt64 `db:"assignee"`
	DueDate         NullTime  `db:"due_date"`
	CompletionTime  NullTime  `db:"completion_time"`
	FeedbackTime    NullTime  `db:"feedback_time"`
	RequestTime     time.Time `db:"request_time"`
	ProgressTime    NullTime  `db:"progress_time"`
}

func (r *DocumentRequest) UnmarshalJSON(data []byte) error {
	return FlexibleJsonStructUnmarshal(data, r)
}
