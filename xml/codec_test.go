package xml_test

import (
	"fmt"

	"github.com/foomo/goencode/xml"
)

func ExampleNewCodec() {
	type Data struct {
		XMLName struct{} `xml:"data"`
		Name    string   `xml:"name"`
	}

	c := xml.NewCodec[Data]()

	encoded, err := c.Encode(Data{Name: "example-123"})
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	fmt.Printf("Encoded: %s\n", string(encoded))
	// Output:
	// Encoded: <data><name>example-123</name></data>
}
