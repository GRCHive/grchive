package core

import (
	"github.com/teambition/rrule-go"
	"time"
)

type Days int32

const (
	KSunday    Days = 0
	KMonday         = 1
	KTuesday        = 2
	KWednesday      = 3
	KThursday       = 4
	KFriday         = 5
	KSaturday       = 6
)

func DaysToRRule(d Days) rrule.Weekday {
	switch d {
	case KSunday:
		return rrule.SU
	case KMonday:
		return rrule.MO
	case KTuesday:
		return rrule.TU
	case KWednesday:
		return rrule.WE
	case KThursday:
		return rrule.TH
	case KFriday:
		return rrule.FR
	case KSaturday:
		return rrule.SA
	}

	// ???
	return rrule.SU
}

func IsPastTime(nowTime time.Time, thresholdTime time.Time, leeway int) bool {
	if nowTime.UTC().Sub(thresholdTime.UTC()).Seconds() > float64(leeway) {
		return true
	}
	return false
}

func CombineDateWithTime(date time.Time, tm time.Time) time.Time {
	// We need this to be in the time's location mainly for RRule parsing
	// since the RRule will take the time of the start time which is what
	// this function is primarily used for atm.
	locDate := date.In(tm.Location())
	locTime := tm

	return time.Date(
		locDate.Year(),
		locDate.Month(),
		locDate.Day(),
		locTime.Hour(),
		locTime.Minute(),
		0,
		0,
		locTime.Location(),
	)
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (r TimeRange) InRange(t time.Time) bool {
	return r.Start.Before(t) && r.End.After(t)
}
