package json

import (
	"encoding/json"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes T to a JSON stream.
func StreamEncoder[T any](w io.Writer, v T) error {
	return json.NewEncoder(w).Encode(v)
}

// StreamDecoder decodes T from a JSON stream.
func StreamDecoder[T any](r io.Reader, v *T) error {
	return json.NewDecoder(r).Decode(v)
}

// NewStreamCodec returns a JSON stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
