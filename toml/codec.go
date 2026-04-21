package toml

import (
	encoding "github.com/foomo/goencode"

	"github.com/BurntSushi/toml"
)

// Encoder encodes T to TOML bytes.
func Encoder[T any](v T) ([]byte, error) {
	return toml.Marshal(v)
}

// Decoder decodes TOML bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return toml.Unmarshal(b, v)
}

// NewCodec returns a TOML codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
