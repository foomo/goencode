package brotli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewEncoder returns a Brotli compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: brotli.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
		buf := sync.Get()
		defer sync.Put(buf)

		w := brotli.NewWriterLevel(buf, o.level)

		if _, err := w.Write(data); err != nil {
			return nil, err
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		return append([]byte(nil), buf.Bytes()...), nil
	}
}

// NewDecoder returns a Brotli decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
		r := brotli.NewReader(bytes.NewReader(data))

		var src io.Reader = r
		if o.maxDecodedSize > 0 {
			src = io.LimitReader(r, o.maxDecodedSize+1)
		}

		decoded, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		if o.maxDecodedSize > 0 && int64(len(decoded)) > o.maxDecodedSize {
			return fmt.Errorf("brotli: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
		}

		*v = decoded

		return nil
	}
}

// NewCodec returns a Brotli compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
