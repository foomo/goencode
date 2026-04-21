package gzip

import (
	"compress/gzip"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a gzip compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	o := options{
		level: gzip.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			gw, err := gzip.NewWriterLevel(w, o.level)
			if err != nil {
				return err
			}

			if _, err := gw.Write(data); err != nil {
				gw.Close()
				return err
			}

			return gw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			gr, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer gr.Close()

			var src io.Reader = gr
			if o.maxDecodedSize > 0 {
				src = io.LimitReader(gr, o.maxDecodedSize+1)
			}

			data, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			if o.maxDecodedSize > 0 && int64(len(data)) > o.maxDecodedSize {
				return fmt.Errorf("gzip: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
			}

			*v = data
			return nil
		},
	}
}
