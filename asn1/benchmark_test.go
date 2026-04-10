package asn1_test

import (
	"testing"

	"github.com/foomo/goencode/asn1"
	"github.com/foomo/goencode/internal/testdata"
)

func BenchmarkCodec(b *testing.B) {
	c := asn1.NewCodec[testdata.User]()

	b.Run("encode", func(b *testing.B) {
		v := *testdata.NewUser()
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode(*testdata.NewUser())
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v testdata.User
			if err := c.Decode(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
