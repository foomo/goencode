package brotli

import (
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	encoding "github.com/foomo/goencode"
)

// NewStreamEncoder returns a Brotli compression stream encoder.
func NewStreamEncoder(opts ...Option) encoding.StreamEncoder[[]byte] {
	o := options{
		level: brotli.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(w io.Writer, data []byte) error {
		bw := brotli.NewWriterLevel(w, o.level)

		if _, err := bw.Write(data); err != nil {
			bw.Close()
			return err
		}

		return bw.Close()
	}
}

// NewStreamDecoder returns a Brotli decompression stream decoder.
func NewStreamDecoder(opts ...Option) encoding.StreamDecoder[[]byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(r io.Reader, v *[]byte) error {
		br := brotli.NewReader(r)

		var src io.Reader = br
		if o.maxDecodedSize > 0 {
			src = io.LimitReader(br, o.maxDecodedSize+1)
		}

		data, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		if o.maxDecodedSize > 0 && int64(len(data)) > o.maxDecodedSize {
			return fmt.Errorf("brotli: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
		}

		*v = data

		return nil
	}
}

// NewStreamCodec returns a Brotli compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: NewStreamEncoder(opts...),
		Decode: NewStreamDecoder(opts...),
	}
}
