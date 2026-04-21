package pem

import (
	stdpem "encoding/pem"
	"errors"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns a PEM codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[*stdpem.Block, []byte] {
	return encoding.Codec[*stdpem.Block, []byte]{
		Encode: func(v *stdpem.Block) ([]byte, error) {
			return stdpem.EncodeToMemory(v), nil
		},
		Decode: func(b []byte, v **stdpem.Block) error {
			block, _ := stdpem.Decode(b)
			if block == nil {
				return errors.New("pem: no PEM block found")
			}

			*v = block

			return nil
		},
	}
}
