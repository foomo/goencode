package flate

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewCodec returns a DEFLATE compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	o := options{
		level: flate.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.Codec[[]byte, []byte]{
		Encode: func(data []byte) ([]byte, error) {
			buf := sync.Get()
			defer sync.Put(buf)

			w, err := flate.NewWriter(buf, o.level)
			if err != nil {
				return nil, err
			}

			if _, err := w.Write(data); err != nil {
				return nil, err
			}

			if err := w.Close(); err != nil {
				return nil, err
			}

			return append([]byte(nil), buf.Bytes()...), nil
		},
		Decode: func(data []byte, v *[]byte) error {
			r := flate.NewReader(bytes.NewReader(data))
			defer r.Close()

			var src io.Reader = r
			if o.maxDecodedSize > 0 {
				src = io.LimitReader(r, o.maxDecodedSize+1)
			}

			decoded, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			if o.maxDecodedSize > 0 && int64(len(decoded)) > o.maxDecodedSize {
				return fmt.Errorf("flate: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
			}

			*v = decoded
			return nil
		},
	}
}
