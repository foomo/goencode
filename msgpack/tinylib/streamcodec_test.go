package msgpack_test

import (
	"bytes"
	"fmt"

	msgpack "github.com/foomo/goencode/msgpack/tinylib"
	"github.com/foomo/goencode/msgpack/tinylib/testdata"
)

func ExampleStreamCodec() {
	c := msgpack.NewStreamCodec[testdata.User]()

	var buf bytes.Buffer
	if err := c.Encode(&buf, *testdata.NewUserTinyLib()); err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	var decoded testdata.User
	if err := c.Decode(&buf, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded Handle: %s\n", decoded.Handle)
	fmt.Printf("Decoded Country: %s\n", decoded.Country)
	// Output:
	// Decoded Handle: @bench
	// Decoded Country: US
}
