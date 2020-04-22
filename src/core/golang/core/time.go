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
	utcDate := date.UTC()
	utcTime := tm.UTC()

	return time.Date(
		utcDate.Year(),
		utcDate.Month(),
		utcDate.Day(),
		utcTime.Hour(),
		utcTime.Minute(),
		0,
		0,
		utcTime.Location(),
	)
}
