package file_test

import (
	"path/filepath"
	"testing"

	"github.com/foomo/goencode/file"
	"github.com/foomo/goencode/internal/testdata"
	json "github.com/foomo/goencode/json/v1"
)

func BenchmarkCodec(b *testing.B) {
	c := file.NewCodec(json.NewCodec[*testdata.User]())
	path := filepath.Join(b.TempDir(), "user.json")

	b.Run("encode", func(b *testing.B) {
		v := testdata.NewUser()
		for b.Loop() {
			if err := c.Encode(path, v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		if err := c.Encode(path, testdata.NewUser()); err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v *testdata.User
			if err := c.Decode(path, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
