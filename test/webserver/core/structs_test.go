package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"testing"
)

func TestStructToMap(t *testing.T) {
	for _, test := range []struct {
		refStruct interface{}
		refMap    map[string]interface{}
	}{
		{
			struct {
				attr1 int
				attr2 string
			}{
				attr1: 32,
				attr2: "bob",
			},
			map[string]interface{}{
				"attr1": 32,
				"attr2": "bob",
			},
		},
		{
			struct {
			}{},
			map[string]interface{}{},
		},
	} {
		testMap := core.StructToMap(test.refStruct)
		assert.Equal(t, len(refMap), len(testMap))
		for rk, rv := range refMap {
			tv, ok := testMap[rk]
			assert.True(t, ok)
			assert.Equal(t, rv, tv)
		}
	}
}
