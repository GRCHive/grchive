package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestSecondsInX(t *testing.T) {
	// 60 seconds in minute, 60 minutes in hour, 24 hours in day
	assert.Equal(t, core.SecondsInDay, 60*60*24)
	assert.Equal(t, core.SecondsInWeek, core.SecondsInDay*7)
	assert.Equal(t, core.SecondsInYear, core.SecondsInDay*365)
}
