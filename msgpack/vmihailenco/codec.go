package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"
)

// Codec is a Codec[T] backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
type Codec[T any] struct{}

// NewCodec returns a msgpack serializer for T.
func NewCodec[T any]() *Codec[T] { return &Codec[T]{} }

func (Codec[T]) Encode(v T) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (Codec[T]) Decode(b []byte, v *T) error {
	return msgpack.Unmarshal(b, v)
}
