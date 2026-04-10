package msgpack

import (
	"fmt"
	"io"

	"github.com/tinylib/msgp/msgp"
)

// StreamCodec is a StreamCodec[T] backed by tinylib/msgp.
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Encodable and msgp.Decodable.
// It is safe for concurrent use.
type StreamCodec[T any] struct{}

// NewStreamCodec returns a msgpack stream serializer for T.
func NewStreamCodec[T any]() *StreamCodec[T] { return &StreamCodec[T]{} }

func (StreamCodec[T]) Encode(w io.Writer, v T) error {
	if e, ok := any(v).(msgp.Encodable); ok {
		return msgp.Encode(w, e)
	}

	if e, ok := any(&v).(msgp.Encodable); ok {
		return msgp.Encode(w, e)
	}

	return fmt.Errorf("msgpack: %T does not implement msgp.Encodable", v)
}

func (StreamCodec[T]) Decode(r io.Reader, v *T) error {
	if d, ok := any(v).(msgp.Decodable); ok {
		return msgp.Decode(r, d)
	}

	return fmt.Errorf("msgpack: %T does not implement msgp.Decodable", v)
}
