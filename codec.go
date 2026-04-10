package goencode

// Codec encodes T to []byte and decodes []byte back to T.
type Codec[T any] interface {
	// Encode encodes v into bytes.
	Encode(v T) ([]byte, error)
	// Decode decodes b into v.
	Decode(b []byte, v *T) error
}
