package main

type Set[T comparable] map[T]bool

func (s Set[T]) Contain(val T) bool {
	_, ok := s[val]
	return ok
}

func (s Set[T]) Insert(val T) {
	s[val] = true
}

func (s Set[T]) Erase(val T) {
	delete(s, val)
}
