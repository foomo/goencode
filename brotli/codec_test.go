package brotli_test

import (
	"fmt"

	goencode "github.com/foomo/goencode"
	"github.com/foomo/goencode/brotli"
	json "github.com/foomo/goencode/json/v1"
)

func ExampleNewCodec() {
	type Data struct {
		Name string
	}

	c := goencode.PipeCodec(json.NewCodec[Data](), brotli.NewCodec())

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
