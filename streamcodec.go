package goencode

import (
	"io"
)

// StreamCodec encodes T to an io.Writer and decodes T from an io.Reader.
type StreamCodec[T any] interface {
	Encode(w io.Writer, v T) error
	Decode(r io.Reader, v *T) error
}
