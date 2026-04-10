package flate

// Option configures a flate Codec.
type Option func(o *options)

type options struct {
	level          int
	maxDecodedSize int64
}

// WithLevel sets the DEFLATE compression level.
// Use compress/flate constants: flate.NoCompression, flate.BestSpeed,
// flate.BestCompression, flate.DefaultCompression, flate.HuffmanOnly.
func WithLevel(level int) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithMaxDecodedSize sets the maximum allowed size of decompressed data in bytes.
// If the decompressed data exceeds this limit, Decode returns an error.
// A value of 0 (the default) means no limit.
func WithMaxDecodedSize(n int64) Option {
	return func(o *options) {
		o.maxDecodedSize = n
	}
}
