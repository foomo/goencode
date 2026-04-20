package goencode

// PipeEncoder chains two encoders: A → B → C.
func PipeEncoder[A, B, C any](first Encoder[A, B], second Encoder[B, C]) Encoder[A, C] {
	return func(a A) (C, error) {
		b, err := first(a)
		if err != nil {
			var zero C
			return zero, err
		}
		return second(b)
	}
}

// PipeDecoder chains two decoders in reverse: decodes C → B via second, then B → A via first.
func PipeDecoder[A, B, C any](first Decoder[A, B], second Decoder[B, C]) Decoder[A, C] {
	return func(c C, a *A) error {
		var b B
		if err := second(c, &b); err != nil {
			return err
		}
		return first(b, a)
	}
}

// PipeCodec chains two codecs: Codec[A,B] + Codec[B,C] → Codec[A,C].
func PipeCodec[A, B, C any](first Codec[A, B], second Codec[B, C]) Codec[A, C] {
	return Codec[A, C]{
		Encode: PipeEncoder(first.Encode, second.Encode),
		Decode: PipeDecoder(first.Decode, second.Decode),
	}
}
