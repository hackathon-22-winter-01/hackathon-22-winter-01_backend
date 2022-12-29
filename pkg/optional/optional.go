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
