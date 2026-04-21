package base64

import (
	stdbase64 "encoding/base64"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes bytes to a Base64 stream.
func StreamEncoder(w io.Writer, v []byte) error {
	enc := stdbase64.NewEncoder(stdbase64.StdEncoding, w)
	if _, err := enc.Write(v); err != nil {
		_ = enc.Close()
		return err
	}

	return enc.Close()
}

// StreamDecoder decodes bytes from a Base64 stream.
func StreamDecoder(r io.Reader, v *[]byte) error {
	data, err := io.ReadAll(stdbase64.NewDecoder(stdbase64.StdEncoding, r))
	if err != nil {
		return err
	}

	*v = data

	return nil
}

// NewStreamCodec returns a Base64 stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
