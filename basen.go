// Package basen implements a radix encoding/decoding scheme, defined by a n-character alphabet.
package basen

import (
	"fmt"
	"math"
	"math/big"
)

// Base58Encoding is the standard base58 encoding, which is the alternate base64 encoding defined in RFC 4648
// modified to avoid both non-alphanumeric characters and letters which might look ambiguous when printed.
// It is designed for human users who manually enter the data, copying from some visual source,
// but also allows easy copy and paste because a double-click will usually select the whole string, and
// it is safe to include in a URL with escaping.
var Base58Encoding = NewEncoding("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base62 is the standard base62 encoding, which is the standard base64 encoding defined in RFC 4648
// modified to avoid non-alphanumeric characters.
// It is useful for generating strings that are safe to include in URLs, and
// allows easy copy and paste because a double-click will usually select the whole string.
var Base62Encoding = NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// An Encoding is a radix encoding/decoding scheme, defined by a n-character alphabet.
type Encoding struct {
	alphabet           []byte
	radix              *big.Int
	bitsPerEncodedByte float64
	decodeMap          [256]int64
}

// NewEncoding returns a new Encoding defined by the given alphabet.
func NewEncoding(alphabet string) *Encoding {
	radix := len(alphabet)

	enc := &Encoding{
		alphabet:           []byte(alphabet),
		radix:              big.NewInt(int64(radix)),
		bitsPerEncodedByte: math.Log2(float64(radix)),
	}

	for i := range enc.decodeMap {
		enc.decodeMap[i] = -1
	}

	for i, c := range alphabet {
		enc.decodeMap[c] = int64(i)
	}

	return enc
}

// EncodedLen returns the length in bytes of the encoding of an input buffer of length n.
func (e *Encoding) EncodedLen(n int) int {
	return int(math.Ceil(float64(n*8) / e.bitsPerEncodedByte))
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of encoded data.
func (e *Encoding) DecodedLen(n int) int {
	return int(math.Ceil(float64(n) * e.bitsPerEncodedByte / 8))
}

// EncodeBigIntToString returns the encoding of src.
func (e *Encoding) EncodeBigIntToString(src *big.Int) string {
	bytes := (src.BitLen() + 7) / 8
	b := make([]byte, 0, e.EncodedLen(bytes))
	rem := new(big.Int)
	zero := new(big.Int)

	for src.Cmp(zero) == 1 {
		src.DivMod(src, e.radix, rem)
		b = append(b, e.alphabet[rem.Int64()])
	}

	reverse(b)

	return string(b)
}

// EncodeInt64ToString returns the encoding of n.
func (e *Encoding) EncodeInt64ToString(n int64) string {
	b := make([]byte, 0, e.EncodedLen(8))
	radix := int64(len(e.alphabet))

	for n > 0 {
		rem := n % radix
		n = n / radix
		b = append(b, e.alphabet[rem])
	}

	reverse(b)

	return string(b)
}

// EncodeToString returns the encoding of src.
func (e *Encoding) EncodeToString(src []byte) string {
	n := new(big.Int)
	n.SetBytes(src)
	return e.EncodeBigIntToString(n)
}

type errInvalidCharacter struct {
	base int64
	char rune
}

func (err errInvalidCharacter) Error() string {
	return fmt.Sprintf("string contains invalid base%d character: %q", err.base, err.char)
}

// DecodeStringToBigInt returns the int64 represented by the encoded string s.
func (e *Encoding) DecodeStringToInt64(s string) (int64, error) {
	var n int64
	radix := e.radix.Int64()

	for _, c := range s {
		i := e.decodeMap[c]
		if i < 0 {
			return 0, errInvalidCharacter{base: e.radix.Int64(), char: c}
		}
		n = n*radix + i
	}

	return n, nil

}

// DecodeStringToBigInt returns the big.Int represented by the encoded string s.
func (e *Encoding) DecodeStringToBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	idx := new(big.Int)

	cap := e.DecodedLen(len(s))
	n.SetBytes(make([]byte, cap))

	for _, c := range s {
		i := e.decodeMap[c]
		if i < 0 {
			return nil, errInvalidCharacter{base: e.radix.Int64(), char: c}
		}
		idx.SetInt64(i)
		n.Add(n.Mul(n, e.radix), idx)
	}

	return n, nil
}

// DecodeString returns the bytes represented by the encoded string s.
func (e *Encoding) DecodeString(s string) ([]byte, error) {
	b, err := e.DecodeStringToBigInt(s)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func reverse(b []byte) {
	for left, right := 0, len(b)-1; left < right; left, right = left+1, right-1 {
		b[left], b[right] = b[right], b[left]
	}
}
