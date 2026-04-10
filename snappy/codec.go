package snappy

import (
	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// Codec is a Codec[T] that applies Snappy compression on top of another Codec[T].
// It is safe for concurrent use.
type Codec[T any] struct {
	codec encoding.Codec[T]
}

// NewCodec returns a Snappy compression codec that delegates serialization to codec.
func NewCodec[T any](codec encoding.Codec[T]) *Codec[T] {
	return &Codec[T]{
		codec: codec,
	}
}

func (c *Codec[T]) Encode(v T) ([]byte, error) {
	b, err := c.codec.Encode(v)
	if err != nil {
		return nil, err
	}

	return snappy.Encode(nil, b), nil
}

func (c *Codec[T]) Decode(b []byte, v *T) error {
	data, err := snappy.Decode(nil, b)
	if err != nil {
		return err
	}

	return c.codec.Decode(data, v)
}
