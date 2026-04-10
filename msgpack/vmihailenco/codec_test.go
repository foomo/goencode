package msgpack_test

import (
	"fmt"

	msgpack "github.com/foomo/goencode/msgpack/vmihailenco"
)

func ExampleCodec() {
	type Data struct {
		Name string `msgpack:"name"`
	}

	c := msgpack.NewCodec[Data]()

	encoded, err := c.Encode(Data{Name: "example-123"})
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	var decoded Data
	if err := c.Decode(encoded, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded Name: %s\n", decoded.Name)
	// Output:
	// Decoded Name: example-123
}
