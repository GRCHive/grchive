package core

import (
	"github.com/teambition/rrule-go"
)

type ModRRuleFn func(*rrule.ROption)

func RebuildRRuleSet(set rrule.Set, modRrule ModRRuleFn) (*rrule.Set, error) {
	newSet := rrule.Set{}

	for _, rr := range set.GetRRule() {
		opts := rr.OrigOptions

		modRrule(&opts)

		newRR, err := rrule.NewRRule(opts)
		if err != nil {
			return nil, err
		}
		newSet.RRule(newRR)
	}
	return &newSet, nil
}
