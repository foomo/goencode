package xml

import (
	"bytes"
	"encoding/xml"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// Encoder encodes T to XML bytes.
func Encoder[T any](v T) ([]byte, error) {
	buf := sync.Get()
	defer sync.Put(buf)

	if err := xml.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

// Decoder decodes XML bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return xml.NewDecoder(bytes.NewReader(b)).Decode(v)
}

// NewCodec returns an XML codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
