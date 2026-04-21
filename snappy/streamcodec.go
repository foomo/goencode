package snappy

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// NewStreamCodec returns a Snappy compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			sw := snappy.NewBufferedWriter(w)
			if _, err := sw.Write(data); err != nil {
				return err
			}
			return sw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			data, err := io.ReadAll(snappy.NewReader(r))
			if err != nil {
				return err
			}
			*v = data
			return nil
		},
	}
}
