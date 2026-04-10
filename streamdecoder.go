package goencode

import (
	"io"
)

type StreamDecoder[T any] interface {
	Decode(r io.Reader, v *T) error
}
