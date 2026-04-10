package gob

import (
	"bytes"
	stdgob "encoding/gob"

	"github.com/foomo/goencode/internal/sync"
)

// Codec is a Codec[T] backed by encoding/gob.
// It is safe for concurrent use.
type Codec[T any] struct{}

// NewCodec returns a GOB serializer for T.
func NewCodec[T any]() *Codec[T] { return &Codec[T]{} }

func (Codec[T]) Encode(v T) ([]byte, error) {
	buf := sync.Get()
	defer sync.Put(buf)

	if err := stdgob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

func (Codec[T]) Decode(b []byte, v *T) error {
	return stdgob.NewDecoder(bytes.NewReader(b)).Decode(v)
}
