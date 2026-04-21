package json

import (
	encoding "github.com/foomo/goencode"
	"github.com/go-json-experiment/json"
)

// Encoder encodes T to JSON bytes (v2).
func Encoder[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

// Decoder decodes JSON bytes into T (v2).
func Decoder[T any](b []byte, v *T) error {
	return json.Unmarshal(b, v)
}

// NewCodec returns a JSON codec for T backed by go-json-experiment/json.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
