package ascii85_test

import (
	"testing"

	"github.com/foomo/goencode/ascii85"
	"github.com/foomo/goencode/internal/testdata"
)

func BenchmarkCodec(b *testing.B) {
	c := ascii85.NewCodec()

	b.Run("encode", func(b *testing.B) {
		v := []byte(testdata.Text)
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode([]byte(testdata.Text))
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v []byte
			if err := c.Decode(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
