package core

import "time"

func IsPastTime(nowTime time.Time, thresholdTime time.Time) bool {
	if nowTime.UTC().Sub(thresholdTime.UTC()).Seconds() > LoadEnvConfig().Login.TimeDriftLeewaySeconds {
		return true
	}
	return false
}
