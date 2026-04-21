package csv_test

import (
	"fmt"

	"github.com/foomo/goencode/csv"
)

func ExampleNewCodec() {
	c := csv.NewCodec()

	records := [][]string{
		{"name", "age"},
		{"Alice", "30"},
	}

	encoded, err := c.Encode(records)
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	fmt.Printf("Encoded: %s", string(encoded))

	var decoded [][]string
	if err := c.Decode(encoded, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded: %v\n", decoded)
	// Output:
	// Encoded: name,age
	// Alice,30
	// Decoded: [[name age] [Alice 30]]
}
