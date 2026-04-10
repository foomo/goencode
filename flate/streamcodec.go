package flate

import (
	"compress/flate"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamCodec is a StreamCodec[T] that applies DEFLATE compression on top of another StreamCodec[T].
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec          encoding.StreamCodec[T]
	level          int
	maxDecodedSize int64
}

// NewStreamCodec returns a flate compression stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		level: flate.DefaultCompression,
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
	fw, err := flate.NewWriter(w, c.level)
	if err != nil {
		return err
	}

	if err := c.codec.Encode(fw, v); err != nil {
		fw.Close()
		return err
	}

	return fw.Close()
}

func (c *StreamCodec[T]) Decode(r io.Reader, v *T) error {
	fr := flate.NewReader(r)
	defer fr.Close()

	var src io.Reader = fr
	if c.maxDecodedSize > 0 {
		src = io.LimitReader(fr, c.maxDecodedSize+1)
	}

	return c.codec.Decode(src, v)
}
