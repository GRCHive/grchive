package core

type Risk struct {
	Id          int64         `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Org         *Organization `db:"org" json:"-"`
}

type RiskFilterData struct {
	NumControls NumericFilterData
}

var NullRiskFilterData RiskFilterData = RiskFilterData{
	NumControls: NullNumericFilterData,
}
