package database

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"time"
)

func buildNumericFilter(ref string, f core.NumericFilterData) string {
	var template string = ""

	switch f.Op {
	case core.Disabled:
		return fmt.Sprintf("%s IS NOT NULL", ref)
	case core.Equal:
		template = "%s = %d"
	case core.NotEqual:
		template = "%s != %d"
	case core.Greater:
		template = "%s > %d"
	case core.GreaterEqual:
		template = "%s >= %d"
	case core.Less:
		template = "%s < %d"
	case core.LessEqual:
		template = "%s <= %d"
	}

	return fmt.Sprintf(template, ref, f.Target)
}

func buildStringFilter(ref string, f core.StringFilterData) string {
	var template string = ""

	switch f.Op {
	case core.SOpDisabled:
		return fmt.Sprintf("%s IS NOT NULL", ref)
	case core.SOpEqual:
		template = "%s = '%s'"
	case core.SOpNotEqual:
		template = "%s != '%s'"
	case core.SOpContains:
		template = "%s LIKE  '%%%s%%'"
	case core.SOpExclude:
		template = "%s NOT LIKE '%%%s%%'"
	}

	return fmt.Sprintf(template, ref, f.Target)
}

func buildTimeRangeFilter(ref string, f core.TimeRangeFilterData) string {
	if !f.Enabled {
		return fmt.Sprintf("%s IS NOT NULL", ref)
	}
	template := "%s >= '%s' AND %s <= '%s'"
	return fmt.Sprintf(template,
		ref, f.Start.Format(time.RFC3339),
		ref, f.End.Format(time.RFC3339))
}
