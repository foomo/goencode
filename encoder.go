package goencode

// Encoder encodes source S to target T.
type Encoder[S, T any] func(s S) (T, error)
