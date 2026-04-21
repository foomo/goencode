package msgpack

import (
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/tinylib/msgp/msgp"
)

// NewStreamCodec returns a msgpack stream codec for T backed by tinylib/msgp.
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Encodable and msgp.Decodable.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			if e, ok := any(v).(msgp.Encodable); ok {
				return msgp.Encode(w, e)
			}
			if e, ok := any(&v).(msgp.Encodable); ok {
				return msgp.Encode(w, e)
			}
			return fmt.Errorf("msgpack: %T does not implement msgp.Encodable", v)
		},
		Decode: func(r io.Reader, v *T) error {
			if d, ok := any(v).(msgp.Decodable); ok {
				return msgp.Decode(r, d)
			}
			return fmt.Errorf("msgpack: %T does not implement msgp.Decodable", v)
		},
	}
}
