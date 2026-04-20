package gob

import (
	"bytes"
	stdgob "encoding/gob"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewCodec returns a GOB codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			buf := sync.Get()
			defer sync.Put(buf)

			if err := stdgob.NewEncoder(buf).Encode(v); err != nil {
				return nil, err
			}

			return append([]byte(nil), buf.Bytes()...), nil
		},
		Decode: func(b []byte, v *T) error {
			return stdgob.NewDecoder(bytes.NewReader(b)).Decode(v)
		},
	}
}
