package goencode

import (
	"io"
)

type StreamEncoder[T any] interface {
	Encode(w io.Writer, v T) error
}
