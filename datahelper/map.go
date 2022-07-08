package datahelper

func MergeMap[K comparable, V any](source, data map[K]V) map[K]V {
	if len(data) == 0 {
		return source
	}
	for k, v := range data {
		source[k] = v
	}
	return source
}
