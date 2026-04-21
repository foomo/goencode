package yaml

import (
	encoding "github.com/foomo/goencode"
	"go.yaml.in/yaml/v4"
)

// Encoder encodes T to YAML v4 bytes.
func Encoder[T any](v T) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decoder decodes YAML v4 bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return yaml.Unmarshal(b, v)
}

// NewCodec returns a YAML v4 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
