package goencode

type Decoder[T any] interface {
	Decode(v any) error
}
