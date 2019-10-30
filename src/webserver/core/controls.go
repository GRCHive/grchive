package core

type ControlType struct {
	Id          int32  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type ControlDocumentationCategory struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ControlId   int64  `db:"control_id"`
}

type Control struct {
	Id                int64     `db:"id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ControlTypeId     int32     `db:"control_type"`
	OrgId             int32     `db:"org_id"`
	FrequencyType     int32     `db:"freq_type"`
	FrequencyInterval int32     `db:"freq_interval"`
	OwnerId           NullInt64 `db:"owner_id"`
}
