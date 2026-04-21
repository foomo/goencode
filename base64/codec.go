package base64

import (
	stdbase64 "encoding/base64"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to Base64.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdbase64.StdEncoding.EncodedLen(len(v)))
	stdbase64.StdEncoding.Encode(dst, v)

	return dst, nil
}

// Decoder decodes Base64 bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdbase64.StdEncoding.DecodedLen(len(b)))

	n, err := stdbase64.StdEncoding.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Base64 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
