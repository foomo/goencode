package file

import "os"

// Option configures a file Codec.
type Option func(o *options)

type options struct {
	perm os.FileMode
}

// WithPermissions sets the file permissions for written files (default: 0o644).
func WithPermissions(perm os.FileMode) Option {
	return func(o *options) {
		o.perm = perm
	}
}
