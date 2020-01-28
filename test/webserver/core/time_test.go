package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
	"time"
)

func TestIsPastTime(t *testing.T) {
	core.InitializeConfig("../../../src/webserver/config/config.toml")
	utcLoc, _ := time.LoadLocation("UTC")
	nyLoc, _ := time.LoadLocation("America/New_York")

	for _, test := range []struct {
		nowTime       time.Time
		thresholdTime time.Time
		isPast        bool
		leeway        int
	}{
		{
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			false,
			1,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, 5, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			false,
			5,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, 8, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			true,
			6,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, nyLoc),
			false,
			5,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, 1, 1, nyLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			true,
			60,
		},
	} {
		testIsPast := core.IsPastTime(test.nowTime, test.thresholdTime, test.leeway)
		assert.Equal(t, test.isPast, testIsPast)
	}
}
