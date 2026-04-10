package goencode

type DecoderFunc[T any] func(v any) error

func (f DecoderFunc[T]) Decode(v any) error {
	return f(v)
}
