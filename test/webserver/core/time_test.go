package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"testing"
	"time"
)

func TestIsPastTime(t *testing.T) {
	utcLoc, _ := time.LoadLocation("UTC")

	var leeway = int(core.LoadEnvConfig().Login.TimeDriftLeewaySeconds)
	for _, test := range []struct {
		nowTime       time.Time
		thresholdTime time.Time
		isPast        bool
	}{
		{
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			false,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, leeway, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			false,
		},
		{
			time.Date(2000, time.January, 1, 1, 1, 2+leeway, 1, utcLoc),
			time.Date(2000, time.January, 1, 1, 1, 1, 1, utcLoc),
			true,
		},
	} {
		testIsPast := core.IsPastTime(test.nowTime, test.thresholdTime)
		assert.Equal(t, test.isPast, testIsPast)
	}
}