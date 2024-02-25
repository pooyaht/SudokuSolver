package solver

import "sync"

type Map[K comparable, V any] interface {
	Store(key K, value V)
	Load(key K) (V, bool)
}

type StandardMap[K comparable, V any] struct {
	data map[K]V
}

func NewStandardMap[K comparable, V any]() StandardMap[K, V] {
	return StandardMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *StandardMap[K, V]) Store(key K, value V) {
	m.data[key] = value
}

func (m StandardMap[K, V]) Load(key K) (value V, ok bool) {
	value, ok = m.data[key]
	return
}

type ConcurrentMap[K comparable, V any] struct {
	data sync.Map
}

func NewConcurrentMap[K comparable, V any]() ConcurrentMap[K, V] {
	return ConcurrentMap[K, V]{}
}

func (m *ConcurrentMap[K, V]) Store(key K, value V) {
	m.data.Store(key, value)
}

func (m *ConcurrentMap[K, V]) Load(key K) (V, bool) {
	val, ok := m.data.Load(key)
	var value V
	if ok {
		value = val.(V)
	}
	return value, ok
}
