package msgpack

import (
	"fmt"

	encoding "github.com/foomo/goencode"
	"github.com/tinylib/msgp/msgp"
)

// NewCodec returns a msgpack codec for T backed by tinylib/msgp.
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Marshaler and msgp.Unmarshaler.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			if m, ok := any(v).(msgp.Marshaler); ok {
				return m.MarshalMsg(nil)
			}
			if m, ok := any(&v).(msgp.Marshaler); ok {
				return m.MarshalMsg(nil)
			}
			return nil, fmt.Errorf("msgpack: %T does not implement msgp.Marshaler", v)
		},
		Decode: func(b []byte, v *T) error {
			if u, ok := any(v).(msgp.Unmarshaler); ok {
				_, err := u.UnmarshalMsg(b)
				return err
			}
			return fmt.Errorf("msgpack: %T does not implement msgp.Unmarshaler", v)
		},
	}
}
