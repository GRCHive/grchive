package core

type Risk struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Identifier  string `db:"identifier"`
	Description string `db:"description"`
	OrgId       int32  `db:"org_id"`
}

type RiskFilterData struct {
	NumControls NumericFilterData
}

var NullRiskFilterData RiskFilterData = RiskFilterData{
	NumControls: NullNumericFilterData,
}
