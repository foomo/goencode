package hex

import (
	stdhex "encoding/hex"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a Hex stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, v []byte) error {
			_, err := stdhex.NewEncoder(w).Write(v)
			return err
		},
		Decode: func(r io.Reader, v *[]byte) error {
			data, err := io.ReadAll(stdhex.NewDecoder(r))
			if err != nil {
				return err
			}
			*v = data
			return nil
		},
	}
}
