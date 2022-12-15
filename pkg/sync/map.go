package sync

import "sync"

// Map sync.Mapの型付きラッパー
type Map[K, V any] struct {
	sm sync.Map
}

func (m *Map[K, V]) Delete(key K) {
	m.sm.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.sm.Load(key)
	if !ok {
		return
	}

	return v.(V), true
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.sm.LoadOrStore(key, value)
	actual = v.(V)

	return
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.sm.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
