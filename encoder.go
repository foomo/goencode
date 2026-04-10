package goencode

type Encoder[T any] interface {
	Encode(v T) error
}
