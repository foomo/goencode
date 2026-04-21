package msgpack

import (
	encoding "github.com/foomo/goencode"
	"github.com/vmihailenco/msgpack/v5"
)

// Encoder encodes T to msgpack bytes (vmihailenco).
func Encoder[T any](v T) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decoder decodes msgpack bytes into T (vmihailenco).
func Decoder[T any](b []byte, v *T) error {
	return msgpack.Unmarshal(b, v)
}

// NewCodec returns a msgpack codec for T backed by vmihailenco/msgpack/v5.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
