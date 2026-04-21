package msgpack

import (
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/vmihailenco/msgpack/v5"
)

// StreamEncoder encodes T to a msgpack stream (vmihailenco).
func StreamEncoder[T any](w io.Writer, v T) error {
	return msgpack.NewEncoder(w).Encode(v)
}

// StreamDecoder decodes T from a msgpack stream (vmihailenco).
func StreamDecoder[T any](r io.Reader, v *T) error {
	return msgpack.NewDecoder(r).Decode(v)
}

// NewStreamCodec returns a msgpack stream codec for T backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: StreamEncoder[T],
		Decode: StreamDecoder[T],
	}
}
