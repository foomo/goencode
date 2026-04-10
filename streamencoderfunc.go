package goencode

import (
	"io"
)

type StreamEncoderFunc[T any] func(w io.Writer, v T) error

func (f StreamEncoderFunc[T]) Encode(w io.Writer, v T) error {
	return f(w, v)
}
