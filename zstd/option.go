package zstd

import "github.com/klauspost/compress/zstd"

// Option configures a zstd Codec.
type Option func(o *options)

type options struct {
	level          zstd.EncoderLevel
	maxDecodedSize int64
}

// WithLevel sets the Zstandard compression level.
// Use zstd constants: zstd.SpeedFastest, zstd.SpeedDefault, zstd.SpeedBetterCompression, zstd.SpeedBestCompression.
func WithLevel(level zstd.EncoderLevel) Option {
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
