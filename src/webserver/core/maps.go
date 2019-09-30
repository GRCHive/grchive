package core

func CopyMap(input map[string]interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})
	for k, v := range input {
		retMap[k] = v
	}
	return retMap
}
