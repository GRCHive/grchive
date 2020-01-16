package core

import "time"

func IsPastTime(nowTime time.Time, thresholdTime time.Time, leeway int) bool {
	if nowTime.UTC().Sub(thresholdTime.UTC()).Seconds() > float64(leeway) {
		return true
	}
	return false
}
