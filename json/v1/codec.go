package json

import (
	"encoding/json"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns a JSON codec for T.
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
