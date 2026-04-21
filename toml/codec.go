package toml

import (
	encoding "github.com/foomo/goencode"

	"github.com/BurntSushi/toml"
)

// NewCodec returns a TOML codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return toml.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			return toml.Unmarshal(b, v)
		},
	}
}
