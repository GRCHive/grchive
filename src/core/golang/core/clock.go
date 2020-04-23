package core

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

var DefaultClock = RealClock{}

type MultiplierClock struct {
	Multiplier int64

	t    time.Time
	lock sync.RWMutex
}

func CreateMultiplierClock(mul int64) *MultiplierClock {
	c := MultiplierClock{
		Multiplier: mul,
		t:          time.Now(),
	}

	go func() {
		for {
			c.lock.Lock()
			c.t = c.t.Add(time.Duration(c.Multiplier) * time.Second)
			c.lock.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

	return &c
}

func (c MultiplierClock) Now() time.Time {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.t
}
