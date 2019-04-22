# basen

[![GoDoc](http://godoc.org/github.com/nicksnyder/basen?status.svg)](http://godoc.org/github.com/nicksnyder/basen) [![Build Status](https://travis-ci.org/nicksnyder/basen.svg?branch=master)](http://travis-ci.org/nicksnyder/basen) [![Report card](https://goreportcard.com/badge/github.com/nicksnyder/basen)](https://goreportcard.com/report/github.com/nicksnyder/basen)

Package basen implements a radix encoding/decoding scheme, defined by a n-character alphabet.

```go
import "github.com/nicksnyder/basen"
```

## base62

Base62 encoding encodes data using the character set `0-9A-Za-z`. It is useful for generating strings that are safe to include in URLs, and allows easy copy and paste because a double-click will usually select the whole string.

Encode:

```go
encoded := basen.Base62Encoding.EncodeToString([]byte("Hello"))
fmt.Println(encoded) // 5TP3P3v
```

Decode:

```go
decoded, _ := basen.Base62Encoding.DecodeString("5TP3P3v")
fmt.Println(string(decoded)) // Hello
```

## base58

[Base58 encoding](https://en.wikipedia.org/wiki/Base58) is base62 encoding modified to avoid letters which might look ambiguous when printed (i.e. `0`, `O`, `I`, `l`). This reduces the risk of errors when a human is manually copying from some visual source.

Encode:

```go
encoded := basen.Base58Encoding.EncodeToString([]byte("Hello"))
fmt.Println(encoded) // 9Ajdvzr
```

Decode:

```go
decoded, _ := basen.Base58Encoding.DecodeString("9Ajdvzr")
fmt.Println(string(decoded)) // Hello
```
