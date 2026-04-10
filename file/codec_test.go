package file_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/foomo/goencode/file"
	json "github.com/foomo/goencode/json/v1"
)

func ExampleCodec() {
	type Data struct {
		Name string
	}

	c := file.NewCodec(json.NewCodec[Data]())

	dir, err := os.MkdirTemp("", "file-codec-example")
	if err != nil {
		fmt.Printf("TempDir failed: %v\n", err)
		return
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "data.json")

	if err := c.Encode(path, Data{Name: "example-123"}); err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	var decoded Data
	if err := c.Decode(path, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded Name: %s\n", decoded.Name)
	// Output:
	// Decoded Name: example-123
}
