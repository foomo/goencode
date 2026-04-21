package ascii85

import (
	"bytes"
	stdascii85 "encoding/ascii85"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to ASCII85.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdascii85.MaxEncodedLen(len(v)))
	n := stdascii85.Encode(dst, v)

	return dst[:n], nil
}

// Decoder decodes ASCII85 bytes.
func Decoder(b []byte, v *[]byte) error {
	buf := bytes.NewBuffer(make([]byte, 0, len(b)))

	r := stdascii85.NewDecoder(bytes.NewReader(b))
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	*v = buf.Bytes()

	return nil
}

// NewCodec returns an ASCII85 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
