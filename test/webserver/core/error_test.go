package core_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"testing"
)

func TestErrorString(t *testing.T) {
	assert.Equal(t, core.ErrorString(nil), "No Error")
	assert.Equal(t, core.ErrorString(errors.New("TEST")), "TEST")
}
