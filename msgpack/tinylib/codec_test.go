package msgpack_test

import (
	"fmt"

	msgpack "github.com/foomo/goencode/msgpack/tinylib"
	"github.com/foomo/goencode/msgpack/tinylib/testdata"
)

func ExampleCodec() {
	c := msgpack.NewCodec[testdata.User]()

	encoded, err := c.Encode(*testdata.NewUserTinyLib())
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	var decoded testdata.User
	if err := c.Decode(encoded, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded Handle: %s\n", decoded.Handle)
	fmt.Printf("Decoded Country: %s\n", decoded.Country)
	// Output:
	// Decoded Handle: @bench
	// Decoded Country: US
}
