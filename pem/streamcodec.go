package pem

import (
	stdpem "encoding/pem"
	"errors"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a PEM stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[*stdpem.Block] {
	return encoding.StreamCodec[*stdpem.Block]{
		Encode: func(w io.Writer, v *stdpem.Block) error {
			return stdpem.Encode(w, v)
		},
		Decode: func(r io.Reader, v **stdpem.Block) error {
			data, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			block, _ := stdpem.Decode(data)
			if block == nil {
				return errors.New("encoding: no PEM block found")
			}
			*v = block
			return nil
		},
	}
}
