package base64

import (
	stdbase64 "encoding/base64"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a Base64 stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, v []byte) error {
			enc := stdbase64.NewEncoder(stdbase64.StdEncoding, w)
			if _, err := enc.Write(v); err != nil {
				_ = enc.Close()
				return err
			}
			return enc.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			data, err := io.ReadAll(stdbase64.NewDecoder(stdbase64.StdEncoding, r))
			if err != nil {
				return err
			}
			*v = data
			return nil
		},
	}
}
