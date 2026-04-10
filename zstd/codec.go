package zstd

import (
	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// Codec is a Codec[T] that applies Zstandard compression on top of another Codec[T].
// It is safe for concurrent use.
type Codec[T any] struct {
	codec          encoding.Codec[T]
	level          zstd.EncoderLevel
	maxDecodedSize int64
}

// NewCodec returns a Zstandard compression codec that delegates serialization to codec.
func NewCodec[T any](codec encoding.Codec[T], opts ...Option) *Codec[T] {
	o := options{
		level: zstd.SpeedDefault,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &Codec[T]{
		codec:          codec,
		level:          o.level,
		maxDecodedSize: o.maxDecodedSize,
	}
}

func (c *Codec[T]) Encode(v T) ([]byte, error) {
	b, err := c.codec.Encode(v)
	if err != nil {
		return nil, err
	}

	enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(c.level))
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	return enc.EncodeAll(b, nil), nil
}

func (c *Codec[T]) Decode(b []byte, v *T) error {
	opts := []zstd.DOption{}
	if c.maxDecodedSize > 0 {
		opts = append(opts, zstd.WithDecoderMaxMemory(uint64(c.maxDecodedSize)))
	}

	dec, err := zstd.NewReader(nil, opts...)
	if err != nil {
		return err
	}
	defer dec.Close()

	data, err := dec.DecodeAll(b, nil)
	if err != nil {
		return err
	}

	return c.codec.Decode(data, v)
}
