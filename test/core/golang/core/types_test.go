package core_test

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"testing"
	"time"
)

type RefTimes struct {
	utcLoc *time.Location
	nyLoc  *time.Location

	refTime  time.Time
	refTime2 time.Time
}

func loadRefTimes() RefTimes {
	ref := RefTimes{}
	ref.utcLoc, _ = time.LoadLocation("UTC")
	ref.nyLoc, _ = time.LoadLocation("America/New_York")

	ref.refTime = time.Date(2000, time.January, 1, 1, 1, 1, 1, ref.utcLoc)
	ref.refTime2 = time.Date(2000, time.January, 1, 1, 1, 1, 1, ref.nyLoc)

	return ref
}
func TestNullTimeMarshalJSON(t *testing.T) {
	refTimes := loadRefTimes()

	for _, ref := range []struct {
		data  core.NullTime
		value interface{}
	}{
		{
			data: core.NullTime{
				sql.NullTime{refTimes.refTime, true}},
			value: refTimes.refTime,
		},
		{
			data: core.NullTime{
				sql.NullTime{refTimes.refTime2, false}},
			value: nil,
		},
	} {
		marshaledJson, _ := json.Marshal(ref.data)

		if ref.value == nil {
			var unmarshalData interface{}
			json.Unmarshal(marshaledJson, &unmarshalData)
			assert.Nil(t, unmarshalData)
		} else {
			var unmarshalData time.Time
			json.Unmarshal(marshaledJson, &unmarshalData)
			assert.Equal(t, ref.value, unmarshalData)
		}
	}
}

func TestNullTimeUnmarshalJSON(t *testing.T) {
	refTimes := loadRefTimes()
	type TestStruct struct {
		Value core.NullTime
	}

	for _, ref := range []struct {
		data  string
		valid bool
		value interface{}
	}{
		{
			data:  "{}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": null}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": \"test\"}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": \"2000-01-01T01:01:01.000000001Z\"}",
			valid: true,
			value: refTimes.refTime,
		},
	} {
		parsed := TestStruct{}
		err := json.Unmarshal([]byte(ref.data), &parsed)
		assert.Nil(t, err)
		assert.Equal(t, ref.valid, parsed.Value.NullTime.Valid)
		if ref.valid {
			assert.Equal(t, ref.value, parsed.Value.NullTime.Time)
		}
	}
}

func TestNullTimeCreate(t *testing.T) {
	refTimes := loadRefTimes()

	for _, ref := range []time.Time{
		refTimes.refTime,
		refTimes.refTime2,
	} {
		test := core.CreateNullTime(ref)
		assert.Equal(t, test.NullTime.Time, ref)
	}
}

func TestNullTimeEqual(t *testing.T) {
	refTimes := loadRefTimes()

	for _, ref := range []struct {
		a     core.NullTime
		b     core.NullTime
		equal bool
	}{
		{
			a:     core.CreateNullTime(refTimes.refTime),
			b:     core.CreateNullTime(refTimes.refTime2),
			equal: false,
		},
		{
			a:     core.CreateNullTime(refTimes.refTime),
			b:     core.CreateNullTime(refTimes.refTime),
			equal: true,
		},
		{
			a:     core.NullTime{sql.NullTime{time.Now(), false}},
			b:     core.CreateNullTime(refTimes.refTime),
			equal: false,
		},
		{
			a:     core.NullTime{sql.NullTime{refTimes.refTime2, false}},
			b:     core.NullTime{sql.NullTime{refTimes.refTime, false}},
			equal: true,
		},
	} {
		assert.Equal(t, ref.equal, ref.a.Equal(ref.b), ref.a.NullTime.Time.String(), ref.b.NullTime.Time.String())
	}
}

func TestNullInt64MarshalJSON(t *testing.T) {
	for _, ref := range []struct {
		data  core.NullInt64
		value interface{}
	}{
		{
			data: core.NullInt64{
				sql.NullInt64{3, true}},
			value: int64(3),
		},
		{
			data: core.NullInt64{
				sql.NullInt64{3, false}},
			value: nil,
		},
	} {
		marshaledJson, _ := json.Marshal(ref.data)

		var unmarshalData interface{}
		json.Unmarshal(marshaledJson, &unmarshalData)

		if ref.value == nil {
			assert.Nil(t, unmarshalData)
		} else {
			assert.Equal(t, ref.value, int64(unmarshalData.(float64)))
		}
	}
}

func TestNullInt64UnmarshalJSON(t *testing.T) {
	type TestStruct struct {
		Value core.NullInt64
	}

	for _, ref := range []struct {
		data  string
		valid bool
		value interface{}
	}{
		{
			data:  "{}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": null}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": \"test\"}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": 5123124}",
			valid: true,
			value: int64(5123124),
		},
	} {
		parsed := TestStruct{}
		err := json.Unmarshal([]byte(ref.data), &parsed)
		assert.Nil(t, err)
		assert.Equal(t, ref.valid, parsed.Value.NullInt64.Valid)
		if ref.valid {
			assert.Equal(t, ref.value, parsed.Value.NullInt64.Int64)
		}
	}
}

func TestNullInt64Create(t *testing.T) {
	for _, ref := range []int64{
		10,
		-245,
		32,
	} {
		test := core.CreateNullInt64(ref)
		assert.Equal(t, test.NullInt64.Int64, ref)
	}
}

func TestNullInt32MarshalJSON(t *testing.T) {
	for _, ref := range []struct {
		data  core.NullInt32
		value interface{}
	}{
		{
			data: core.NullInt32{
				sql.NullInt32{3, true}},
			value: int32(3),
		},
		{
			data: core.NullInt32{
				sql.NullInt32{3, false}},
			value: nil,
		},
	} {
		marshaledJson, _ := json.Marshal(ref.data)

		var unmarshalData interface{}
		json.Unmarshal(marshaledJson, &unmarshalData)

		if ref.value == nil {
			assert.Nil(t, unmarshalData)
		} else {
			assert.Equal(t, ref.value, int32(unmarshalData.(float64)))
		}
	}
}

func TestNullInt32UnmarshalJSON(t *testing.T) {
	type TestStruct struct {
		Value core.NullInt32
	}

	for _, ref := range []struct {
		data  string
		valid bool
		value interface{}
	}{
		{
			data:  "{}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": null}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": \"test\"}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": 5123124}",
			valid: true,
			value: int32(5123124),
		},
	} {
		parsed := TestStruct{}
		err := json.Unmarshal([]byte(ref.data), &parsed)
		assert.Nil(t, err)
		assert.Equal(t, ref.valid, parsed.Value.NullInt32.Valid)
		if ref.valid {
			assert.Equal(t, ref.value, parsed.Value.NullInt32.Int32)
		}
	}
}

func TestNullInt32Create(t *testing.T) {
	for _, ref := range []int32{
		10,
		-245,
		32,
	} {
		test := core.CreateNullInt32(ref)
		assert.Equal(t, test.NullInt32.Int32, ref)
	}
}

func TestNullBoolMarshalJSON(t *testing.T) {
	for _, ref := range []struct {
		data  core.NullBool
		value interface{}
	}{
		{
			data: core.NullBool{
				sql.NullBool{true, true}},
			value: true,
		},
		{
			data: core.NullBool{
				sql.NullBool{false, true}},
			value: false,
		},
		{
			data: core.NullBool{
				sql.NullBool{false, false}},
			value: nil,
		},
	} {
		marshaledJson, _ := json.Marshal(ref.data)

		var unmarshalData interface{}
		json.Unmarshal(marshaledJson, &unmarshalData)

		if ref.value == nil {
			assert.Nil(t, unmarshalData)
		} else {
			assert.Equal(t, ref.value, unmarshalData.(bool))
		}
	}
}

func TestNullBoolUnmarshalJSON(t *testing.T) {
	type TestStruct struct {
		Value core.NullBool
	}

	for _, ref := range []struct {
		data  string
		valid bool
		value interface{}
	}{
		{
			data:  "{}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": null}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": \"test\"}",
			valid: false,
			value: nil,
		},
		{
			data:  "{\"Value\": true}",
			valid: true,
			value: true,
		},
		{
			data:  "{\"Value\": false}",
			valid: true,
			value: false,
		},
	} {
		parsed := TestStruct{}
		err := json.Unmarshal([]byte(ref.data), &parsed)
		assert.Nil(t, err)
		assert.Equal(t, ref.valid, parsed.Value.NullBool.Valid)
		if ref.valid {
			assert.Equal(t, ref.value, parsed.Value.NullBool.Bool)
		}
	}
}

func TestNullBoolCreate(t *testing.T) {
	for _, ref := range []bool{
		true,
		false,
	} {
		test := core.CreateNullBool(ref)
		assert.Equal(t, test.NullBool.Bool, ref)
	}
}
