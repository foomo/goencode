package csv

import (
	"bytes"
	stdcsv "encoding/csv"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// Encoder encodes [][]string to CSV bytes.
func Encoder(v [][]string) ([]byte, error) {
	buf := sync.Get()
	defer sync.Put(buf)

	cw := stdcsv.NewWriter(buf)
	if err := cw.WriteAll(v); err != nil {
		return nil, err
	}

	cw.Flush()

	if err := cw.Error(); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

// Decoder decodes CSV bytes into [][]string.
func Decoder(b []byte, v *[][]string) error {
	records, err := stdcsv.NewReader(bytes.NewReader(b)).ReadAll()
	if err != nil {
		return err
	}

	*v = records

	return nil
}

// NewCodec returns a CSV codec for [][]string.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[][]string, []byte] {
	return encoding.Codec[[][]string, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
