package goencode

import "io"

// StreamEncoder encodes S into an io.Writer.
type StreamEncoder[S any] func(w io.Writer, s S) error

// StreamDecoder decodes S from an io.Reader.
type StreamDecoder[S any] func(r io.Reader, s *S) error

// StreamCodec bundles streaming encode/decode for S.
type StreamCodec[S any] struct {
	Encode StreamEncoder[S]
	Decode StreamDecoder[S]
}
