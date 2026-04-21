package yaml

import (
	"io"

	encoding "github.com/foomo/goencode"
	"go.yaml.in/yaml/v3"
)

// StreamEncoder encodes T to a YAML v3 stream.
func StreamEncoder[T any](w io.Writer, v T) error {
	return yaml.NewEncoder(w).Encode(v)
}

// StreamDecoder decodes T from a YAML v3 stream.
func StreamDecoder[T any](r io.Reader, v *T) error {
	return yaml.NewDecoder(r).Decode(v)
}

// NewStreamCodec returns a YAML v3 stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
