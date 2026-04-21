package asn1

import (
	stdasn1 "encoding/asn1"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes T to ASN.1 bytes.
func Encoder[T any](v T) ([]byte, error) {
	return stdasn1.Marshal(v)
}

// Decoder decodes ASN.1 bytes into T.
func Decoder[T any](b []byte, v *T) error {
	_, err := stdasn1.Unmarshal(b, v)
	return err
}

// NewCodec returns an ASN1 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
