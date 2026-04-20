package base32_test

import (
	"fmt"

	"github.com/foomo/goencode/base32"
)

func ExampleNewCodec() {
	c := base32.NewCodec()

	encoded, err := c.Encode([]byte("hello"))
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	fmt.Printf("Encoded: %s\n", string(encoded))
	// Output:
	// Encoded: NBSWY3DP
}
