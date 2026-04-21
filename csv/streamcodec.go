package csv

import (
	stdcsv "encoding/csv"
	"io"

	encoding "github.com/foomo/goencode"
)

// StreamEncoder encodes [][]string to a CSV stream.
func StreamEncoder(w io.Writer, v [][]string) error {
	cw := stdcsv.NewWriter(w)
	if err := cw.WriteAll(v); err != nil {
		return err
	}

	cw.Flush()

	return cw.Error()
}

// StreamDecoder decodes [][]string from a CSV stream.
func StreamDecoder(r io.Reader, v *[][]string) error {
	records, err := stdcsv.NewReader(r).ReadAll()
	if err != nil {
		return err
	}

	*v = records

	return nil
}

// NewStreamCodec returns a CSV stream codec for [][]string.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[][]string] {
	return encoding.StreamCodec[[][]string]{
		Encode: StreamEncoder,
		Decode: StreamDecoder,
	}
}
