package msgpack

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
)

// Codec is a Codec[T] backed by tinylib/msgp.
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Marshaler and msgp.Unmarshaler.
// It is safe for concurrent use.
type Codec[T any] struct{}

// NewCodec returns a msgpack serializer for T.
func NewCodec[T any]() *Codec[T] { return &Codec[T]{} }

func (Codec[T]) Encode(v T) ([]byte, error) {
	if m, ok := any(v).(msgp.Marshaler); ok {
		return m.MarshalMsg(nil)
	}

	if m, ok := any(&v).(msgp.Marshaler); ok {
		return m.MarshalMsg(nil)
	}

	return nil, fmt.Errorf("msgpack: %T does not implement msgp.Marshaler", v)
}

func (Codec[T]) Decode(b []byte, v *T) error {
	if u, ok := any(v).(msgp.Unmarshaler); ok {
		_, err := u.UnmarshalMsg(b)
		return err
	}

	return fmt.Errorf("msgpack: %T does not implement msgp.Unmarshaler", v)
}
