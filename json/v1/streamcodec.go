package json

import (
	"encoding/json"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a JSON stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return json.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v *T) error {
			return json.NewDecoder(r).Decode(v)
		},
	}
}
