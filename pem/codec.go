package pem

import (
	stdpem "encoding/pem"
	"errors"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes a PEM block to bytes.
func Encoder(v *stdpem.Block) ([]byte, error) {
	return stdpem.EncodeToMemory(v), nil
}

// Decoder decodes bytes into a PEM block.
func Decoder(b []byte, v **stdpem.Block) error {
	block, _ := stdpem.Decode(b)
	if block == nil {
		return errors.New("pem: no PEM block found")
	}

	*v = block

	return nil
}

// NewCodec returns a PEM codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[*stdpem.Block, []byte] {
	return encoding.Codec[*stdpem.Block, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
