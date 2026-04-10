package file

import (
	"fmt"
	"os"
	"path/filepath"

	encoding "github.com/foomo/goencode"
)

// StreamCodec encodes T to a file and decodes T from a file using an underlying StreamCodec[T].
// Writes are atomic: data is written to a temporary file and renamed into place.
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec encoding.StreamCodec[T]
	perm  os.FileMode
}

// NewStreamCodec returns a file stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		perm: 0o644,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &StreamCodec[T]{
		codec: codec,
		perm:  o.perm,
	}
}

// Encode serializes v and atomically writes the result to path.
func (c *StreamCodec[T]) Encode(path string, v T) error {
	dir := filepath.Dir(path)

	f, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	tmp := f.Name()

	if err := c.codec.Encode(f, v); err != nil {
		f.Close()
		os.Remove(tmp)

		return fmt.Errorf("encoding to temp file: %w", err)
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
func (c *StreamCodec[T]) Decode(path string, v *T) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.codec.Decode(f, v)
}
