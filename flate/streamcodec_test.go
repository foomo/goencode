package flate_test

import (
	"bytes"
	"fmt"

	"github.com/foomo/goencode/flate"
)

func ExampleNewStreamCodec() {
	c := flate.NewStreamCodec()

	input := []byte("hello flate stream")
	var buf bytes.Buffer
	if err := c.Encode(&buf, input); err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	var decoded []byte
	if err := c.Decode(&buf, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded: %s\n", string(decoded))
	// Output:
	// Decoded: hello flate stream
}
