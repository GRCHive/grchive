package core

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
