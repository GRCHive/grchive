package core

const SearchNotFound int = -1

func LinearSearchInt32Slice(arr []int32, v int32) int {
	for i, t := range arr {
		if t == v {
			return i
		}
	}
	return SearchNotFound
}
