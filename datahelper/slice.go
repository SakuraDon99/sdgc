package datahelper

func NewSlice[T any](size ...int) []T {
	s := 0
	if len(size) > 0 {
		s = size[len(size)-1]
	}
	return make([]T, s)
}

func SliceContains[T comparable](arr []T, val T) bool {
	for _, t := range arr {
		if val == t {
			return true
		}
	}
	return false
}

func SliceContainsF[T any](arr []T, val T, f func(src, dst T) bool) bool {
	for _, t := range arr {
		if f(t, val) {
			return true
		}
	}
	return false
}

func SliceGetOrDefault[T any](arr []T, i int, def T) T {
	if len(arr)-1 < i {
		return def
	}
	return arr[i]
}
