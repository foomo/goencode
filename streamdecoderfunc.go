package goencode

import (
	"io"
)

type StreamDecoderFunc[T any] func(r io.Reader, v *T) error

func (f StreamDecoderFunc[T]) Decode(r io.Reader, v *T) error {
	return f(r, v)
}
