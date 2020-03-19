package core

func CreateMapFromKeyValues(keys []string, values []interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})
	for idx, key := range keys {
		if idx >= len(values) {
			break
		}
		retMap[key] = values[idx]
	}
	return retMap
}

func CopyMap(input map[string]interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})
	for k, v := range input {
		retMap[k] = v
	}
	return retMap
}

func MergeMaps(input ...map[string]interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})
	for i := 0; i < len(input); i++ {
		for k, v := range input[i] {
			retMap[k] = v
		}
	}
	return retMap
}

func GetFromMapWithKeyOptions(input map[string]interface{}, keys ...string) interface{} {
	for _, k := range keys {
		if k == "" {
			continue
		}

		v, ok := input[k]
		if !ok {
			continue
		}
		return v
	}
	return nil
}
