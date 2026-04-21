package yaml

import (
	encoding "github.com/foomo/goencode"
	"go.yaml.in/yaml/v3"
)

// NewCodec returns a YAML v3 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return yaml.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			return yaml.Unmarshal(b, v)
		},
	}
}
