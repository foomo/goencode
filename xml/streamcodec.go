package xml

import (
	"encoding/xml"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes T to an XML stream.
func StreamEncoder[T any](w io.Writer, v T) error {
	return xml.NewEncoder(w).Encode(v)
}

// StreamDecoder decodes T from an XML stream.
func StreamDecoder[T any](r io.Reader, v *T) error {
	return xml.NewDecoder(r).Decode(v)
}

// NewStreamCodec returns an XML stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
