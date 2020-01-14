package core

import (
	"fmt"
	"time"
)

type timer struct {
	startTime time.Time
	lastTime  time.Time
}

func NewTimer() timer {
	currentTime := time.Now()
	return timer{
		startTime: currentTime,
		lastTime:  currentTime,
	}
}

func (t timer) LogTick(s string) {
	currentTime := time.Now()
	total := currentTime.Sub(t.startTime)
	elapsed := currentTime.Sub(t.lastTime)

	Info(fmt.Sprintf("%s :: Total (%dms) :: Elapsed (%dms)", s, total.Milliseconds(), elapsed.Milliseconds()))
}
