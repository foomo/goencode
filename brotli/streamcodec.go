package brotli

import (
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a Brotli compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	o := options{
		level: brotli.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			bw := brotli.NewWriterLevel(w, o.level)

			if _, err := bw.Write(data); err != nil {
				bw.Close()
				return err
			}

			return bw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
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
		},
	}
}
