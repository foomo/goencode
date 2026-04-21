package msgpack

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/vmihailenco/msgpack/v5"
)

// NewStreamCodec returns a msgpack stream codec for T backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return msgpack.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v *T) error {
			return msgpack.NewDecoder(r).Decode(v)
		},
	}
}
