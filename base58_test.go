package basen_test

import (
	"fmt"

	"github.com/nicksnyder/basen"
)

func ExampleBase58Encoding() {
	encoded := basen.Base58Encoding.EncodeToString([]byte("Hello"))
	fmt.Println(encoded)

	decoded, err := basen.Base58Encoding.DecodeString(encoded)
	fmt.Println(string(decoded), err)

	// Output:
	// 9Ajdvzr
	// Hello <nil>
}
