package asn1

import (
	stdasn1 "encoding/asn1"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes T to an ASN.1 stream.
func StreamEncoder[T any](w io.Writer, v T) error {
	data, err := stdasn1.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}

// StreamDecoder decodes T from an ASN.1 stream.
func StreamDecoder[T any](r io.Reader, v *T) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	_, err = stdasn1.Unmarshal(data, v)

	return err
}

// NewStreamCodec returns an ASN.1 stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
