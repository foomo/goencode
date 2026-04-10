package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// StreamCodec is a StreamCodec[T] backed by github.com/BurntSushi/toml.
// It is safe for concurrent use.
type StreamCodec[T any] struct{}

// NewStreamCodec returns a TOML stream codec for T.
func NewStreamCodec[T any]() *StreamCodec[T] { return &StreamCodec[T]{} }

func (StreamCodec[T]) Encode(w io.Writer, v T) error {
	return toml.NewEncoder(w).Encode(v)
}

func (StreamCodec[T]) Decode(r io.Reader, v *T) error {
	_, err := toml.NewDecoder(r).Decode(v)
	return err
}
