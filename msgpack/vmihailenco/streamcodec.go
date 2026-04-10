package msgpack

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// StreamCodec is a StreamCodec[T] backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
type StreamCodec[T any] struct{}

// NewStreamCodec returns a msgpack stream serializer for T.
func NewStreamCodec[T any]() *StreamCodec[T] { return &StreamCodec[T]{} }

func (StreamCodec[T]) Encode(w io.Writer, v T) error {
	return msgpack.NewEncoder(w).Encode(v)
}

func (StreamCodec[T]) Decode(r io.Reader, v *T) error {
	return msgpack.NewDecoder(r).Decode(v)
}
