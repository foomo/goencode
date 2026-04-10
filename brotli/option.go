package brotli

// Option configures a brotli Codec.
type Option func(o *options)

type options struct {
	level          int
	maxDecodedSize int64
}

// WithLevel sets the Brotli compression level (0–11).
// Higher values yield better compression at the cost of speed.
// Use brotli.DefaultCompression (6), brotli.BestSpeed (0), or brotli.BestCompression (11).
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
