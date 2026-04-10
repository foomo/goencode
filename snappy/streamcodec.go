package snappy

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// StreamCodec is a StreamCodec[T] that applies Snappy compression on top of another StreamCodec[T].
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec encoding.StreamCodec[T]
}

// NewStreamCodec returns a Snappy compression stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T]) *StreamCodec[T] {
	return &StreamCodec[T]{
		codec: codec,
	}
}

func (c *StreamCodec[T]) Encode(w io.Writer, v T) error {
	sw := snappy.NewBufferedWriter(w)

	if err := c.codec.Encode(sw, v); err != nil {
		return err
	}

	return sw.Close()
}

func (c *StreamCodec[T]) Decode(r io.Reader, v *T) error {
	return c.codec.Decode(snappy.NewReader(r), v)
}
