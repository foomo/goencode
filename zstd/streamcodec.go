package zstd

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// NewStreamCodec returns a Zstandard compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	o := options{
		level: zstd.SpeedDefault,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			zw, err := zstd.NewWriter(w, zstd.WithEncoderLevel(o.level))
			if err != nil {
				return err
			}

			if _, err := zw.Write(data); err != nil {
				zw.Close()
				return err
			}

			return zw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			dopts := []zstd.DOption{}
			if o.maxDecodedSize > 0 {
				dopts = append(dopts, zstd.WithDecoderMaxMemory(uint64(o.maxDecodedSize)))
			}

			zr, err := zstd.NewReader(r, dopts...)
			if err != nil {
				return err
			}
			defer zr.Close()

			data, err := io.ReadAll(zr)
			if err != nil {
				return err
			}

			*v = data
			return nil
		},
	}
}
