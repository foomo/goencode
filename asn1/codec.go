package asn1

import (
	stdasn1 "encoding/asn1"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns an ASN1 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return stdasn1.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			_, err := stdasn1.Unmarshal(b, v)
			return err
		},
	}
}
