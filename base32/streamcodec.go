package base32

import (
	stdbase32 "encoding/base32"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes bytes to a Base32 stream.
func StreamEncoder(w io.Writer, v []byte) error {
	enc := stdbase32.NewEncoder(stdbase32.StdEncoding, w)
	if _, err := enc.Write(v); err != nil {
		return err
	}

	return enc.Close()
}

// StreamDecoder decodes bytes from a Base32 stream.
func StreamDecoder(r io.Reader, v *[]byte) error {
	data, err := io.ReadAll(stdbase32.NewDecoder(stdbase32.StdEncoding, r))
	if err != nil {
		return err
	}

	*v = data

	return nil
}

// NewStreamCodec returns a Base32 stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
