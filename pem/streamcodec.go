package pem

import (
	stdpem "encoding/pem"
	"errors"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes a PEM block to a stream.
func StreamEncoder(w io.Writer, v *stdpem.Block) error {
	return stdpem.Encode(w, v)
}

// StreamDecoder decodes a PEM block from a stream.
func StreamDecoder(r io.Reader, v **stdpem.Block) error {
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
}

// NewStreamCodec returns a PEM stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[*stdpem.Block] {
	return encoding.StreamCodec[*stdpem.Block]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
