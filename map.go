package main

type Map[K comparable, V interface{}] map[K]V

func (m Map[K, V]) Contain(key K) bool {
	_, ok := m[key]
	return ok
}
