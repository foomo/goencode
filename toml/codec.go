package toml

import (
	"github.com/BurntSushi/toml"
)

// Codec is a Codec[T] backed by github.com/BurntSushi/toml.
// It is safe for concurrent use.
type Codec[T any] struct{}

// NewCodec returns a TOML codec for T.
func NewCodec[T any]() Codec[T] { return Codec[T]{} }

func (Codec[T]) Encode(v T) ([]byte, error) {
	return toml.Marshal(v)
}

func (Codec[T]) Decode(b []byte, v *T) error {
	return toml.Unmarshal(b, v)
}
