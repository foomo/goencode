package flate

import (
	"compress/flate"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a DEFLATE compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	o := options{
		level: flate.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			fw, err := flate.NewWriter(w, o.level)
			if err != nil {
				return err
			}

			if _, err := fw.Write(data); err != nil {
				fw.Close()
				return err
			}

			return fw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			fr := flate.NewReader(r)
			defer fr.Close()

			var src io.Reader = fr
			if o.maxDecodedSize > 0 {
				src = io.LimitReader(fr, o.maxDecodedSize+1)
			}

			data, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			if o.maxDecodedSize > 0 && int64(len(data)) > o.maxDecodedSize {
				return fmt.Errorf("flate: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
			}

			*v = data

			return nil
		},
	}
}
