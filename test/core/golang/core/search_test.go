package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestLinearSearchInt32Slice(t *testing.T) {
	for _, ref := range []struct {
		slice  []int32
		search int32
		index  int
	}{
		{
			slice:  []int32{5, 3, 2, 10},
			search: 2,
			index:  2,
		},
		{
			slice:  []int32{5, 3, 2, 10},
			search: 5,
			index:  0,
		},
		{
			slice:  []int32{5, 3, 2, 10},
			search: 10,
			index:  3,
		},
		{
			slice:  []int32{5, 3, 2, 10},
			search: 22,
			index:  core.SearchNotFound,
		},
	} {
		assert.Equal(t, ref.index, core.LinearSearchInt32Slice(ref.slice, ref.search))
	}
}
