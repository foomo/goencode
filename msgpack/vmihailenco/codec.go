package msgpack

import (
	encoding "github.com/foomo/goencode"
	"github.com/vmihailenco/msgpack/v5"
)

// NewCodec returns a msgpack codec for T backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}
