package base32

import (
	stdbase32 "encoding/base32"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to Base32.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdbase32.StdEncoding.EncodedLen(len(v)))
	stdbase32.StdEncoding.Encode(dst, v)

	return dst, nil
}

// Decoder decodes Base32 bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdbase32.StdEncoding.DecodedLen(len(b)))

	n, err := stdbase32.StdEncoding.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Base32 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
