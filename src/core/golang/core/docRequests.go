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
	ApproveTime     NullTime  `db:"approve_time"`
	ApproveUserId   NullInt64 `db:"approve_user_id"`
}

func (r *DocumentRequest) UnmarshalJSON(data []byte) error {
	return FlexibleJsonStructUnmarshal(data, r)
}

type DocRequestStatus int32

const (
	DocRequestOpen       DocRequestStatus = 0
	DocRequestInProgress                  = 1
	DocRequestFeedback                    = 2
	DocRequestComplete                    = 3
	DocRequestOverdue                     = 4
	DocRequestApproved                    = 5
)

type DocRequestStatusFilterData struct {
	ValidStatuses []DocRequestStatus
}

type DocRequestFilterData struct {
	RequestTimeFilter TimeRangeFilterData
	DueDateFilter     TimeRangeFilterData
	StatusFilter      DocRequestStatusFilterData
	RequesterFilter   UserFilterData
	AssigneeFilter    UserFilterData
}

var ValidDueDateDocRequestFilter = DocRequestFilterData{
	DueDateFilter: TimeRangeFilterData{
		Enabled: true,
		// Just an arbitrarily large date
		End: CreateNullTime(time.Date(3000, 12, 1, 0, 0, 0, 0, time.UTC)),
	},
}
