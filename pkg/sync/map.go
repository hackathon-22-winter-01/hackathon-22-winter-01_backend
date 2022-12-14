package sync

import "sync"

// Map is a wrapper of sync.Map to provide type-safe methods.
type Map[K, V any] struct {
	sync.Map
}

func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.Map.LoadOrStore(key, value)
	actual = v.(V)
	return
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
