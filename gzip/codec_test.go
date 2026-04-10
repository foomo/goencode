package gzip_test

import (
	"fmt"

	"github.com/foomo/goencode/gzip"
	json "github.com/foomo/goencode/json/v1"
)

func ExampleCodec() {
	type Data struct {
		Name string
	}

	c := gzip.NewCodec(json.NewCodec[Data]())

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
