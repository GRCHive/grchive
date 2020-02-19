package database

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
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
