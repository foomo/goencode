package msgpack_test

import (
	"testing"

	"github.com/foomo/goencode/internal/testdata"
	msgpack "github.com/foomo/goencode/msgpack/vmihailenco"
)

func BenchmarkCodec(b *testing.B) {
	c := msgpack.NewCodec[*testdata.User]()

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
