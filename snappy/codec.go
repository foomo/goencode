package snappy

import (
	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// Encoder compresses bytes using Snappy.
func Encoder(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

// Decoder decompresses Snappy bytes.
func Decoder(data []byte, v *[]byte) error {
	decoded, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	*v = decoded

	return nil
}

// NewCodec returns a Snappy compression codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
