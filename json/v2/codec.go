package json

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/go-json-experiment/json"
)

// NewCodec returns a JSON codec for T backed by go-json-experiment/json.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return json.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			return json.Unmarshal(b, v)
		},
	}
}

// NewStreamCodec returns a JSON stream codec for T backed by go-json-experiment/json.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return json.MarshalWrite(w, v)
		},
		Decode: func(r io.Reader, v *T) error {
			return json.UnmarshalRead(r, v)
		},
	}
}
