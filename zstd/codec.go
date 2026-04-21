package zstd

import (
	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// NewCodec returns a Zstandard compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	o := options{
		level: zstd.SpeedDefault,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.Codec[[]byte, []byte]{
		Encode: func(data []byte) ([]byte, error) {
			enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(o.level))
			if err != nil {
				return nil, err
			}
			defer enc.Close()

			return enc.EncodeAll(data, nil), nil
		},
		Decode: func(data []byte, v *[]byte) error {
			dopts := []zstd.DOption{}
			if o.maxDecodedSize > 0 {
				dopts = append(dopts, zstd.WithDecoderMaxMemory(uint64(o.maxDecodedSize)))
			}

			dec, err := zstd.NewReader(nil, dopts...)
			if err != nil {
				return err
			}
			defer dec.Close()

			decoded, err := dec.DecodeAll(data, nil)
			if err != nil {
				return err
			}

			*v = decoded

			return nil
		},
	}
}
