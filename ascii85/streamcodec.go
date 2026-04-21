package ascii85

import (
	stdascii85 "encoding/ascii85"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes bytes to an ASCII85 stream.
func StreamEncoder(w io.Writer, v []byte) error {
	dst := make([]byte, stdascii85.MaxEncodedLen(len(v)))
	n := stdascii85.Encode(dst, v)
	_, err := w.Write(dst[:n])

	return err
}

// StreamDecoder decodes bytes from an ASCII85 stream.
func StreamDecoder(r io.Reader, v *[]byte) error {
	data, err := io.ReadAll(stdascii85.NewDecoder(r))
	if err != nil {
		return err
	}

	*v = data

	return nil
}

// NewStreamCodec returns an ASCII85 stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
