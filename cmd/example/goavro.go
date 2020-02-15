package main

import (
	"fmt"

	"github.com/linkedin/goavro"
)

func main() {
	codec, err := goavro.NewCodec(`
            {
              "type": "record",
              "name": "LongList",
              "fields" : [
	      {"name": "next", "type": ["null", "LongList", {"type": "long", "logicalType": "timestamp-millis"}], "default": null}
              ]
            }`)
	if err != nil {
		fmt.Println(err)
	}

	// NOTE: May omit fields when using default value
	textual := []byte(`{"next":{"LongList":{}}}`)

	// Convert textual Avro data (in Avro JSON format) to native Go form
	native, _, err := codec.NativeFromTextual(textual)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("nativeFromTextual %v\n", native)

	// Convert native Go form to binary Avro data
	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("BinaryFromNative %v\n", binary)

	// Convert binary Avro data back to native Go form
	native, _, err = codec.NativeFromBinary(binary)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("NativeFromBinary %v\n", native)

	// Convert native Go form to textual Avro data
	textual, err = codec.TextualFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("TextualFromNative %v\n", binary)
}
