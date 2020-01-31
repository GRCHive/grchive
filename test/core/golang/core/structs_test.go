package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
)

func TestStructToMap(t *testing.T) {
	for _, test := range []struct {
		refStruct interface{}
		refMap    map[string]interface{}
	}{
		{
			struct {
				Attr1 int
				Attr2 string
			}{
				Attr1: 32,
				Attr2: "bob",
			},
			map[string]interface{}{
				"Attr1": 32,
				"Attr2": "bob",
			},
		},
		{
			struct {
			}{},
			map[string]interface{}{},
		},
	} {
		testMap := core.StructToMap(test.refStruct)
		assert.Equal(t, len(test.refMap), len(testMap))
		for rk, rv := range test.refMap {
			tv, ok := testMap[rk]
			assert.True(t, ok)
			assert.Equal(t, rv, tv)
		}
	}
}
