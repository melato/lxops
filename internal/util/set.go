package util

type Set[T comparable] map[T]struct{}

func (s Set[T]) Put(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}
