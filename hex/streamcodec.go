package hex

import (
	stdhex "encoding/hex"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes bytes to a hexadecimal stream.
func StreamEncoder(w io.Writer, v []byte) error {
	_, err := stdhex.NewEncoder(w).Write(v)
	return err
}

// StreamDecoder decodes bytes from a hexadecimal stream.
func StreamDecoder(r io.Reader, v *[]byte) error {
	data, err := io.ReadAll(stdhex.NewDecoder(r))
	if err != nil {
		return err
	}

	*v = data

	return nil
}

// NewStreamCodec returns a Hex stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
