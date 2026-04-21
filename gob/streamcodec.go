package gob

import (
	"encoding/gob"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes T to a gob stream.
func StreamEncoder[T any](w io.Writer, v T) error {
	return gob.NewEncoder(w).Encode(v)
}

// StreamDecoder decodes T from a gob stream.
func StreamDecoder[T any](r io.Reader, v *T) error {
	return gob.NewDecoder(r).Decode(v)
}

// NewStreamCodec returns a GOB stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
