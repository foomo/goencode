package base64_test

import (
	"fmt"

	"github.com/foomo/goencode/base64"
)

func ExampleNewCodec() {
	c := base64.NewCodec()

	encoded, err := c.Encode([]byte("hello"))
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	fmt.Printf("Encoded: %s\n", string(encoded))
	// Output:
	// Encoded: aGVsbG8=
}
