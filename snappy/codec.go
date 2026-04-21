package snappy

import (
	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// NewCodec returns a Snappy compression codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: func(data []byte) ([]byte, error) {
			return snappy.Encode(nil, data), nil
		},
		Decode: func(data []byte, v *[]byte) error {
			decoded, err := snappy.Decode(nil, data)
			if err != nil {
				return err
			}
			*v = decoded
			return nil
		},
	}
}
