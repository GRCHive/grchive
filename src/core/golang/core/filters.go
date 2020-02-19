package core

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
