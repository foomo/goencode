package snappy

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// StreamEncoder compresses bytes to a Snappy stream.
func StreamEncoder(w io.Writer, data []byte) error {
	sw := snappy.NewBufferedWriter(w)
	if _, err := sw.Write(data); err != nil {
		return err
	}

	return sw.Close()
}

// StreamDecoder decompresses bytes from a Snappy stream.
func StreamDecoder(r io.Reader, v *[]byte) error {
	data, err := io.ReadAll(snappy.NewReader(r))
	if err != nil {
		return err
	}

	*v = data

	return nil
}

// NewStreamCodec returns a Snappy compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
