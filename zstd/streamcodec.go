package zstd

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// StreamCodec is a StreamCodec[T] that applies Zstandard compression on top of another StreamCodec[T].
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec          encoding.StreamCodec[T]
	level          zstd.EncoderLevel
	maxDecodedSize int64
}

// NewStreamCodec returns a Zstandard compression stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		level: zstd.SpeedDefault,
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
	zw, err := zstd.NewWriter(w, zstd.WithEncoderLevel(c.level))
	if err != nil {
		return err
	}

	if err := c.codec.Encode(zw, v); err != nil {
		zw.Close()
		return err
	}

	return zw.Close()
}

func (c *StreamCodec[T]) Decode(r io.Reader, v *T) error {
	opts := []zstd.DOption{}
	if c.maxDecodedSize > 0 {
		opts = append(opts, zstd.WithDecoderMaxMemory(uint64(c.maxDecodedSize)))
	}

	zr, err := zstd.NewReader(r, opts...)
	if err != nil {
		return err
	}
	defer zr.Close()

	return c.codec.Decode(zr, v)
}
