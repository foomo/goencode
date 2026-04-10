package file

import (
	"fmt"
	"os"
	"path/filepath"

	encoding "github.com/foomo/goencode"
)

// Codec encodes T to a file and decodes T from a file using an underlying Codec[T].
// Writes are atomic: data is written to a temporary file and renamed into place.
// It is safe for concurrent use.
type Codec[T any] struct {
	codec encoding.Codec[T]
	perm  os.FileMode
}

// NewCodec returns a file codec that delegates serialization to codec.
func NewCodec[T any](codec encoding.Codec[T], opts ...Option) *Codec[T] {
	o := options{
		perm: 0o644,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &Codec[T]{
		codec: codec,
		perm:  o.perm,
	}
}

// Encode serializes v and atomically writes the result to path.
func (c *Codec[T]) Encode(path string, v T) error {
	b, err := c.codec.Encode(v)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	f, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	tmp := f.Name()

	if _, err := f.Write(b); err != nil {
		f.Close()
		os.Remove(tmp)

		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Chmod(tmp, c.perm); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("setting file permissions: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// Decode reads the file at path and deserializes its contents into v.
func (c *Codec[T]) Decode(path string, v *T) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return c.codec.Decode(b, v)
}
