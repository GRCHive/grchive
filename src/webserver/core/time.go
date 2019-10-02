package core

import "time"

func IsPastTime(inTime time.Time) bool {
	if time.Now().UTC().Sub(inTime).Seconds() > LoadEnvConfig().Login.TimeDriftLeewaySeconds {
		return true
	}
	return false
}
