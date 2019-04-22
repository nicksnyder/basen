package basen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"
)

var tests = []struct {
	encoding  *Encoding
	decoded   string
	encoded   string
	decodeErr error
}{
	{
		encoding: Base62Encoding,
		decoded:  " ",
		encoded:  "W",
	},
	{
		encoding: Base58Encoding,
		decoded:  " ",
		encoded:  "Z",
	},
	{
		encoding: Base62Encoding,
		decoded:  "hello world",
		encoded:  "AAwf93rvy4aWQVw",
	},
	{
		encoding: Base58Encoding,
		decoded:  "hello world",
		encoded:  "StV1DL6CwTryKyV",
	},
	{
		encoding:  Base62Encoding,
		encoded:   "-",
		decodeErr: errInvalidCharacter{base: 62, char: '-'},
	},
	{
		encoding:  Base58Encoding,
		encoded:   "-",
		decodeErr: errInvalidCharacter{base: 58, char: '-'},
	},
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		if test.decodeErr != nil {
			continue
		}

		t.Run(test.decoded, func(t *testing.T) {
			actual := test.encoding.EncodeToString([]byte(test.decoded))
			if actual != test.encoded {
				t.Fatalf("got %q; want %q", actual, test.encoded)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for _, test := range tests {
		t.Run(test.decoded, func(t *testing.T) {
			actual, err := test.encoding.DecodeString(test.encoded)
			if err != test.decodeErr {
				t.Fatalf("got error %#v; want %#v", err, test.decodeErr)
			}
			if string(actual) != test.decoded {
				t.Fatalf("got %q; want %q", actual, test.encoded)
			}
		})
	}
}

func TestBase62EncodeDecode(t *testing.T) {
	bytez := []byte{0x01, 0xaa, 0xff}
	for _, b := range bytez {
		expected := []byte{b}
		t.Run(fmt.Sprintf("%#x", b), func(t *testing.T) {
			for i := 0; i < 5; i++ {
				expected = append(expected, expected...)
				t.Run(fmt.Sprintf("%d", 2<<uint(1)), func(t *testing.T) {
					s := Base62Encoding.EncodeToString(expected)
					actual, err := Base62Encoding.DecodeString(s)
					if err != nil {
						t.Fatalf("error decoding %q: %s", s, err)
					}
					if !bytes.Equal(actual, expected) {
						t.Fatalf("%q got %x; want %x", s, actual, expected)
					}
				})
			}
		})
	}
}

func TestBase58EncodeDecode(t *testing.T) {
	bytez := []byte{0x01, 0xaa, 0xff}
	for _, b := range bytez {
		expected := []byte{b}
		t.Run(fmt.Sprintf("%#x", b), func(t *testing.T) {
			for i := 0; i < 5; i++ {
				expected = append(expected, expected...)
				t.Run(fmt.Sprintf("%d", 2<<uint(1)), func(t *testing.T) {
					s := Base58Encoding.EncodeToString(expected)
					actual, err := Base58Encoding.DecodeString(s)
					if err != nil {
						t.Fatalf("error decoding %q: %s", s, err)
					}
					if !bytes.Equal(actual, expected) {
						t.Fatalf("got %x; want %x", actual, expected)
					}
				})
			}
		})
	}
}

var benchmarkEncodeResult string

func BenchmarkBase62EncodeToString(b *testing.B) {
	buf := []byte{0xff}
	for i := 0; i < 4; i++ {
		buf = append(buf, buf...)
		b.Run(fmt.Sprintf("%d bytes", 2<<uint(i)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkEncodeResult = Base62Encoding.EncodeToString(buf)
			}
		})
	}
}

var benchmarkDecodeResult []byte

func BenchmarkBase62DecodeString(b *testing.B) {
	buf := []byte{0xff}
	for i := 0; i < 4; i++ {
		buf = append(buf, buf...)
		enc := Base62Encoding.EncodeToString(buf)
		b.Run(fmt.Sprintf("%d bytes", 2<<uint(i)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkDecodeResult, _ = Base62Encoding.DecodeString(enc)
			}
		})
	}
}

func BenchmarkBase64EncodeToString(b *testing.B) {
	buf := []byte{0xff}
	for i := 0; i < 4; i++ {
		buf = append(buf, buf...)
		b.Run(fmt.Sprintf("%d bytes", 2<<uint(i)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkEncodeResult = base64.StdEncoding.EncodeToString(buf)
			}
		})
	}
}

func BenchmarkBase62EncodeBigIntToString(b *testing.B) {
	buf := []byte{0xff}
	for i := 0; i < 4; i++ {
		buf = append(buf, buf...)
		n := new(big.Int).SetBytes(buf)
		b.Run(fmt.Sprintf("%d bytes", 2<<uint(i)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkEncodeResult = Base62Encoding.EncodeBigIntToString(n)
			}
		})
	}
}
