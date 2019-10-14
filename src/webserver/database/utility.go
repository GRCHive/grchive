package database

import (
	"encoding/json"
	"github.com/lib/pq"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

// Checks whether the error indicates a duplicate entry on INSERT.
func IsDuplicateDBEntry(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case *pq.Error:
		return err.(*pq.Error).Code == "23505"
	}
	return false
}

func readProcessFlowInputOutputArray(data []uint8) ([]core.ProcessFlowInputOutput, error) {
	// Manually do the unmarshal so we don't have to add JSON field tags to the
	// ProcessFlowInputOutput structure because I don't want to have those fields
	// be named differently when sent via JSON.
	mapArr := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(data), &mapArr)
	if err != nil {
		return nil, err
	}

	retArr := make([]core.ProcessFlowInputOutput, len(mapArr))
	for i := 0; i < len(mapArr); i++ {
		retArr[i].Id = int64(mapArr[i]["id"].(float64))
		retArr[i].Name = mapArr[i]["name"].(string)
		retArr[i].ParentNodeId = int64(mapArr[i]["parent_node_id"].(float64))
		retArr[i].TypeId = int32(mapArr[i]["io_type_id"].(float64))
	}
	return retArr, nil
}
