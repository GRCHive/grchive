package core_test

import (
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"testing"
)

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
