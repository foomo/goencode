package hex

import (
	stdhex "encoding/hex"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to hexadecimal.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdhex.EncodedLen(len(v)))
	stdhex.Encode(dst, v)

	return dst, nil
}

// Decoder decodes hexadecimal bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdhex.DecodedLen(len(b)))

	n, err := stdhex.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Hex codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
