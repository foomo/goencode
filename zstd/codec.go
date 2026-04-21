package zstd

import (
	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// NewEncoder returns a Zstandard compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: zstd.SpeedDefault,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
		enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(o.level))
		if err != nil {
			return nil, err
		}
		defer enc.Close()

		return enc.EncodeAll(data, nil), nil
	}
}

// NewDecoder returns a Zstandard decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
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
	}
}

// NewCodec returns a Zstandard compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
