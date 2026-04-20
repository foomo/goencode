package goencode

// Encoder encodes source S to target T.
type Encoder[S, T any] func(s S) (T, error)

// Decoder decodes target T back into source S.
type Decoder[S, T any] func(t T, s *S) error

// Codec bundles an Encoder and Decoder for S ↔ T round-trips.
type Codec[S, T any] struct {
	Encode Encoder[S, T]
	Decode Decoder[S, T]
}
