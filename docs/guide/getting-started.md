# Getting Started

## Installation

```bash
go get github.com/foomo/goencode
```

Some packages live in their own Go modules and must be installed separately:

::: warning Submodule Packages
These packages have external dependencies and require their own `go get`:

| Package | Install |
|---------|---------|
| `json/v2` | `go get github.com/foomo/goencode/json/v2` |
| `yaml/v2` | `go get github.com/foomo/goencode/yaml/v2` |
| `yaml/v3` | `go get github.com/foomo/goencode/yaml/v3` |
| `yaml/v4` | `go get github.com/foomo/goencode/yaml/v4` |
| `toml` | `go get github.com/foomo/goencode/toml` |
| `snappy` | `go get github.com/foomo/goencode/snappy` |
| `zstd` | `go get github.com/foomo/goencode/zstd` |
| `brotli` | `go get github.com/foomo/goencode/brotli` |
| `msgpack/tinylib` | `go get github.com/foomo/goencode/msgpack/tinylib` |
| `msgpack/vmihailenco` | `go get github.com/foomo/goencode/msgpack/vmihailenco` |
:::

## Core Types

goencode defines function types and struct bundles at the root package.

### Codec[S, T] — byte-oriented

```go
// Function types
type Encoder[S, T any] func(s S) (T, error)
type Decoder[S, T any] func(t T, s *S) error

// Codec bundles an Encoder and Decoder pair.
type Codec[S, T any] struct {
    Encode Encoder[S, T]
    Decode Decoder[S, T]
}
```

Use `Codec[S, T]` when you need the encoded result as a value — for example, `Codec[User, []byte]` for serialization or `Codec[[]byte, []byte]` for compression. Codecs compose via `PipeCodec` for type-safe chaining.

### StreamCodec[S] — io.Reader/io.Writer-oriented

```go
// Stream function types
type StreamEncoder[S any] func(w io.Writer, s S) error
type StreamDecoder[S any] func(r io.Reader, s *S) error

// StreamCodec bundles a StreamEncoder and StreamDecoder pair.
type StreamCodec[S any] struct {
    Encode StreamEncoder[S]
    Decode StreamDecoder[S]
}
```

Use `StreamCodec[S]` when working with streams — HTTP request/response bodies, files, network connections, or any `io.Reader`/`io.Writer`.

## Minimal Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/json/v1"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    c := json.NewCodec[User]() // [!code highlight]

    // Encode
    b, err := c.Encode(User{Name: "Alice", Age: 30}) // [!code highlight]
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(b)) // {"name":"Alice","age":30}

    // Decode
    var u User
    if err := c.Decode(b, &u); err != nil { // [!code highlight]
        log.Fatal(err)
    }
    fmt.Println(u.Name) // Alice
}
```

## Concurrency Safety

All codecs in this library are safe for concurrent use. Codec structs bundle pure function values with no shared mutable state. You can safely share a single codec instance across goroutines.
