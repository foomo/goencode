package goencode

type EncoderFunc[T any] func(v T) error

func (f EncoderFunc[T]) Encode(v T) error {
	return f(v)
}
