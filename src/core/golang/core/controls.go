package core

type ControlType struct {
	Id          int32  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type Control struct {
	Id                int64     `db:"id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ControlTypeId     int32     `db:"control_type"`
	OrgId             int32     `db:"org_id"`
	FrequencyType     int32     `db:"freq_type"`
	FrequencyInterval int32     `db:"freq_interval"`
	FrequencyOther    string    `db:"freq_other"`
	OwnerId           NullInt64 `db:"owner_id"`
	Manual            bool      `db:"is_manual"`
}

type ControlFilterData struct {
	NumRisks NumericFilterData
}

var NullControlFilterData ControlFilterData = ControlFilterData{
	NumRisks: NullNumericFilterData,
}
