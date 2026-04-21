package toml

import (
	"io"

	encoding "github.com/foomo/goencode"

	"github.com/BurntSushi/toml"
)

// NewStreamCodec returns a TOML stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return toml.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v *T) error {
			_, err := toml.NewDecoder(r).Decode(v)

			return err
		},
	}
}
