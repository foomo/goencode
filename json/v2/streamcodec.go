package json

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/go-json-experiment/json"
)

// StreamEncoder encodes T to a JSON stream (v2).
func StreamEncoder[T any](w io.Writer, v T) error {
	return json.MarshalWrite(w, v)
}

// StreamDecoder decodes T from a JSON stream (v2).
func StreamDecoder[T any](r io.Reader, v *T) error {
	return json.UnmarshalRead(r, v)
}

// NewStreamCodec returns a JSON stream codec for T backed by go-json-experiment/json.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
