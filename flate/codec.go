package flate

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// Codec is a Codec[T] that applies DEFLATE compression on top of another Codec[T].
// It is safe for concurrent use.
type Codec[T any] struct {
	codec          encoding.Codec[T]
	level          int
	maxDecodedSize int64
}

// NewCodec returns a flate compression codec that delegates serialization to codec.
func NewCodec[T any](codec encoding.Codec[T], opts ...Option) *Codec[T] {
	o := options{
		level: flate.DefaultCompression,
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

	buf := sync.Get()
	defer sync.Put(buf)

	w, err := flate.NewWriter(buf, c.level)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(b); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

func (c *Codec[T]) Decode(b []byte, v *T) error {
	r := flate.NewReader(bytes.NewReader(b))
	defer r.Close()

	var src io.Reader = r
	if c.maxDecodedSize > 0 {
		src = io.LimitReader(r, c.maxDecodedSize+1)
	}

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	if c.maxDecodedSize > 0 && int64(len(data)) > c.maxDecodedSize {
		return fmt.Errorf("flate: decompressed size exceeds limit of %d bytes", c.maxDecodedSize)
	}

	return c.codec.Decode(data, v)
}
