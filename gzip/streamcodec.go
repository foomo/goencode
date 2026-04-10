package gzip

import (
	"compress/gzip"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamCodec is a StreamCodec[T] that applies gzip compression on top of another StreamCodec[T].
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec          encoding.StreamCodec[T]
	level          int
	maxDecodedSize int64
}

// NewStreamCodec returns a gzip compression stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		level: gzip.DefaultCompression,
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
	gw, err := gzip.NewWriterLevel(w, c.level)
	if err != nil {
		return err
	}

	if err := c.codec.Encode(gw, v); err != nil {
		gw.Close()
		return err
	}

	return gw.Close()
}

func (c *StreamCodec[T]) Decode(r io.Reader, v *T) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	var src io.Reader = gr
	if c.maxDecodedSize > 0 {
		src = io.LimitReader(gr, c.maxDecodedSize+1)
	}

	return c.codec.Decode(src, v)
}
