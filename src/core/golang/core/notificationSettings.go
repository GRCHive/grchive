package core

type PbcNotificationCadenceSettings struct {
	Id              int64   `db:"id"`
	OrgId           int32   `db:"org_id"`
	DaysBeforeDue   int32   `db:"days_before_due"`
	SendToAssignee  bool    `db:"send_to_assignee"`
	SendToRequester bool    `db:"send_to_requester"`
	AdditionalUsers []int64 `db:"additional_users"`
}
