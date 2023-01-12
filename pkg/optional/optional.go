package optional

type Of[T any] struct {
	V  T
	OK bool
}

func NewFromPtr[T any](v *T) Of[T] {
	if v == nil {
		return Of[T]{}
	}

	return Of[T]{V: *v, OK: true}
}

func (v Of[T]) Value() (T, bool) {
	return v.V, v.OK
}
