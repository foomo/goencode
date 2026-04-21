package gzip_test

import (
	"bytes"
	"fmt"

	"github.com/foomo/goencode/gzip"
)

func ExampleNewStreamCodec() {
	c := gzip.NewStreamCodec()

	input := []byte("hello gzip stream")

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
	// Decoded: hello gzip stream
}
