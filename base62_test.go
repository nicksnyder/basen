package basen_test

import (
	"fmt"

	"github.com/nicksnyder/basen"
)

func ExampleBase62Encoding() {
	encoded := basen.Base62Encoding.EncodeToString([]byte("Hello"))
	fmt.Println(encoded)

	decoded, err := basen.Base62Encoding.DecodeString(encoded)
	fmt.Println(string(decoded), err)

	// Output:
	// 5TP3P3v
	// Hello <nil>
}
