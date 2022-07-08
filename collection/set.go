package collection

type Set[T any] struct {
	m map[any]bool
}

func NewSet[T any]() *Set[T] {
	return &Set[T]{
		m: make(map[any]bool),
	}
}

func (s *Set[T]) Put(val T) {
	s.m[val] = true
}

func (s *Set[T]) Contains(val T) bool {
	_, ok := s.m[val]
	return ok
}

func (s *Set[T]) ToSlice() []T {
	var arr []T
	for val := range s.m {
		arr = append(arr, val.(T))
	}
	return arr
}

func (s *Set[T]) Len() int {
	return len(s.m)
}
