package csv_test

import (
	"testing"

	"github.com/foomo/goencode/csv"
)

func BenchmarkCodec(b *testing.B) {
	c := csv.NewCodec()

	b.Run("encode", func(b *testing.B) {
		v := [][]string{
			{"name", "age", "country"},
			{"Alice", "30", "US"},
			{"Bob", "25", "UK"},
			{"Charlie", "35", "DE"},
		}
		for b.Loop() {
			if _, err := c.Encode(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("decode", func(b *testing.B) {
		data, err := c.Encode([][]string{
			{"name", "age", "country"},
			{"Alice", "30", "US"},
			{"Bob", "25", "UK"},
			{"Charlie", "35", "DE"},
		})
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var v [][]string
			if err := c.Decode(data, &v); err != nil {
				b.Fatal(err)
			}
		}
	})
}
