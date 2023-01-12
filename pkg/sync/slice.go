package sync

import "sync"

type Slice[T any] struct {
	s   []T
	mux sync.Mutex
}

func NewSlice[T any]() Slice[T] {
	return Slice[T]{
		s:   make([]T, 0),
		mux: sync.Mutex{},
	}
}

func (s *Slice[T]) Append(value T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s = append(s.s, value)
}

func (s *Slice[T]) Clone() []T {
	s.mux.Lock()
	defer s.mux.Unlock()

	return append([]T{}, s.s...)
}
