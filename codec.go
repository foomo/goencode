package goencode

// Codec bundles an Encoder and Decoder for S ↔ T round-trips.
type Codec[S, T any] struct {
	Encode Encoder[S, T]
	Decode Decoder[S, T]
}
