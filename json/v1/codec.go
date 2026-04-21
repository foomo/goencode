package json

import (
	"encoding/json"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes T to JSON bytes (v1).
func Encoder[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

// Decoder decodes JSON bytes into T (v1).
func Decoder[T any](b []byte, v *T) error {
	return json.Unmarshal(b, v)
}

// NewCodec returns a JSON codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
