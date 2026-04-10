package pem_test

import (
	stdpem "encoding/pem"
	"testing"

	"github.com/foomo/goencode/internal/testdata"
	"github.com/foomo/goencode/pem"
)

func BenchmarkCodec(b *testing.B) {
	c := pem.NewCodec()

	b.Run("encode", func(b *testing.B) {
		v := &stdpem.Block{
			Type:  "TEST",
			Bytes: []byte(testdata.Text),
		}
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode(&stdpem.Block{
			Type:  "TEST",
			Bytes: []byte(testdata.Text),
		})
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v *stdpem.Block
			if err := c.Decode(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
