package core

import (
	"time"
)

type ComparisonOperators int

const (
	Disabled ComparisonOperators = iota
	Equal
	NotEqual
	Greater
	GreaterEqual
	Less
	LessEqual
)

type NumericFilterData struct {
	Op     ComparisonOperators
	Target int
}

var NullNumericFilterData NumericFilterData = NumericFilterData{
	Op: Disabled,
}

type StringComparisonOperators int

const (
	SOpDisabled StringComparisonOperators = iota
	SOpEqual
	SOpNotEqual
	SOpContains
	SOpExclude
)

type StringFilterData struct {
	Op     StringComparisonOperators
	Target string
}

var NullStringFilterData StringFilterData = StringFilterData{
	Op: SOpDisabled,
}

type TimeRangeFilterData struct {
	Enabled bool
	Start   time.Time
	End     time.Time
}
