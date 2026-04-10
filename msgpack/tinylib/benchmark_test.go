package msgpack_test

import (
	"testing"

	msgpack "github.com/foomo/goencode/msgpack/tinylib"
	"github.com/foomo/goencode/msgpack/tinylib/testdata"
)

func BenchmarkCodec(b *testing.B) {
	c := msgpack.NewCodec[testdata.User]()

	b.Run("encode", func(b *testing.B) {
		v := *testdata.NewUserTinyLib()
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode(*testdata.NewUserTinyLib())
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
