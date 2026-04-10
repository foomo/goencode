package yaml_test

import (
	"testing"

	"github.com/foomo/goencode/internal/testdata"
	"github.com/foomo/goencode/yaml/v2"
)

func BenchmarkCodec(b *testing.B) {
	c := yaml.NewCodec[*testdata.User]()

	b.Run("encode", func(b *testing.B) {
		v := testdata.NewUser()
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode(testdata.NewUser())
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v *testdata.User
			if err := c.Decode(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
