package core_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"reflect"
	"testing"
)

var refKeys []string = []string{
	"KEY1",
	"kEy2",
	".",
}

var refValues []interface{} = []interface{}{
	124354,
	nil,
	"hello",
}

var refKeys2 []string = []string{
	"ALT-KEY1",
	"ALT-kEy2",
	"ALT-.",
}

var refValues2 []interface{} = []interface{}{
	54321,
	"test",
	nil,
}

func TestCreateMapFromKeyValues(t *testing.T) {
	testMap := core.CreateMapFromKeyValues(refKeys, refValues)
	for idx, key := range refKeys {
		val, ok := testMap[key]
		assert.True(t, ok)
		assert.Equal(t, val, refValues[idx])
	}
}

func TestCreateMapFromKeyValuesUnequalArrayLength(t *testing.T) {
	// # Keys < # Values, create a valid map with available keys (ignore extra values)
	testMap := core.CreateMapFromKeyValues([]string{"hello"}, refValues)
	assert.True(t, reflect.DeepEqual(testMap, map[string]interface{}{"hello": refValues[0]}))

	// # Keys > # Values, create a valid map with available values (ignore extra keys)
	testMap = core.CreateMapFromKeyValues(refKeys, []interface{}{"hello"})
	assert.True(t, reflect.DeepEqual(testMap, map[string]interface{}{refKeys[0]: "hello"}))
}

func TestCopyMap(t *testing.T) {
	refMap := core.CreateMapFromKeyValues(refKeys, refValues)
	testMap := core.CopyMap(refMap)

	// Ensure the maps are equal immediately after copy.
	assert.True(t, reflect.DeepEqual(refMap, testMap))

	// Ensure the maps don't point to the same object.
	assert.True(t, &refMap != &testMap)

	// Ensure that modifying the refMap has no effect on the test map.
	for _, key := range refKeys {
		oldVal := refMap[key]
		refMap[key] = 99999999
		assert.Equal(t, oldVal, testMap[key])
	}
}

func TestMergeMapEmpty(t *testing.T) {
	// Merging nothing should result in nothing.
	newMap := core.MergeMaps()
	assert.True(t, reflect.DeepEqual(newMap, make(map[string]interface{})))
}

func TestMergeMapsSingle(t *testing.T) {
	// Merging a single map is equivalent to a copy.
	testMap1 := core.CreateMapFromKeyValues(refKeys, refValues)
	newMap := core.MergeMaps(testMap1)
	assert.True(t, reflect.DeepEqual(testMap1, newMap))
}

func TestMergeMapsNoOverlapTwo(t *testing.T) {
	// Merging multiple maps with no overlaps creates a map with both sets of keys.
	testMap1 := core.CreateMapFromKeyValues(refKeys, refValues)
	testMap2 := core.CreateMapFromKeyValues(refKeys2, refValues2)
	newMap := core.MergeMaps(testMap1, testMap2)
	for idx, key := range refKeys {
		val, ok := newMap[key]
		assert.True(t, ok)
		assert.Equal(t, val, refValues[idx])
	}

	for idx, key := range refKeys2 {
		val, ok := newMap[key]
		assert.True(t, ok)
		assert.Equal(t, val, refValues2[idx])
	}
}

func TestMergeMapsOverlap(t *testing.T) {
	// Merging multiple maps with overlaps will take the values of the later map.
	testMap1 := core.CreateMapFromKeyValues(refKeys, refValues)
	testMap2 := core.CreateMapFromKeyValues(refKeys, refValues2)
	newMap := core.MergeMaps(testMap1, testMap2)
	assert.True(t, reflect.DeepEqual(testMap2, newMap))

	newMap = core.MergeMaps(testMap2, testMap1)
	assert.True(t, reflect.DeepEqual(testMap1, newMap))
}
