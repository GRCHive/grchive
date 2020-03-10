package core

import "time"

type AuditEvent struct {
	Id                int64                  `db:"id"`
	OrgId             int32                  `db:"org_id"`
	ResourceType      string                 `db:"resource_type"`
	ResourceId        string                 `db:"resource_id"`
	ResourceExtraData map[string]interface{} `db:"resource_extra_data"`
	Action            string                 `db:"action"`
	PerformedAt       time.Time              `db:"performed_at"`
	UserId            NullInt64              `db:"user_id"`
}

type AuditTrailRetrievalParams struct {
	OrgId        NullInt32
	EventId      NullInt64
	ResourceType []string
	ResourceId   []string
}

type AuditTrailSortParams struct {
	SortColumns   []string
	SortDirection NullString
	Limit         NullInt32
	Page          NullInt32
}

type AuditTrailFilterData struct {
	ResourceTypeFilter StringFilterData
	ActionFilter       StringFilterData
	UserFilter         StringFilterData
	TimeRangeFilter    TimeRangeFilterData
}

var NullAuditTrailFilterData AuditTrailFilterData = AuditTrailFilterData{}
