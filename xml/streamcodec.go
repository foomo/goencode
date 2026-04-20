package xml

import (
	"encoding/xml"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns an XML stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return xml.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v *T) error {
			return xml.NewDecoder(r).Decode(v)
		},
	}
}
