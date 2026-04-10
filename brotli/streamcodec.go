package brotli

import (
	"io"

	"github.com/andybalholm/brotli"
	encoding "github.com/foomo/goencode"
)

// StreamCodec is a StreamCodec[T] that applies Brotli compression on top of another StreamCodec[T].
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec          encoding.StreamCodec[T]
	level          int
	maxDecodedSize int64
}

// NewStreamCodec returns a Brotli compression stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		level: brotli.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &StreamCodec[T]{
		codec:          codec,
		level:          o.level,
		maxDecodedSize: o.maxDecodedSize,
	}
}

func (c *StreamCodec[T]) Encode(w io.Writer, v T) error {
	bw := brotli.NewWriterLevel(w, c.level)

	if err := c.codec.Encode(bw, v); err != nil {
		bw.Close()
		return err
	}

	return bw.Close()
}

func (c *StreamCodec[T]) Decode(r io.Reader, v *T) error {
	br := brotli.NewReader(r)

	var src io.Reader = br
	if c.maxDecodedSize > 0 {
		src = io.LimitReader(br, c.maxDecodedSize+1)
	}

	return c.codec.Decode(src, v)
}
