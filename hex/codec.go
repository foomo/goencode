package hex

import (
	stdhex "encoding/hex"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns a Hex codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: func(v []byte) ([]byte, error) {
			dst := make([]byte, stdhex.EncodedLen(len(v)))
			stdhex.Encode(dst, v)

			return dst, nil
		},
		Decode: func(b []byte, v *[]byte) error {
			dst := make([]byte, stdhex.DecodedLen(len(b)))

			n, err := stdhex.Decode(dst, b)
			if err != nil {
				return err
			}

			*v = dst[:n]

			return nil
		},
	}
}
