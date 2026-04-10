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
| `json2` | `go get github.com/foomo/goencode/json2` |
| `yaml/v2` | `go get github.com/foomo/goencode/yaml/v2` |
| `yaml/v3` | `go get github.com/foomo/goencode/yaml/v3` |
| `yaml/v4` | `go get github.com/foomo/goencode/yaml/v4` |
| `snappy` | `go get github.com/foomo/goencode/snappy` |
| `zstd` | `go get github.com/foomo/goencode/zstd` |
:::

## Core Interfaces

goencode defines two generic interfaces at the root package.

### Codec[T] — byte-oriented

```go
// Codec encodes T to []byte and decodes []byte back to T.
type Codec[T any] interface {
    Encode(v T) ([]byte, error)
    Decode(b []byte, v *T) error
}
```

Use `Codec[T]` when you need the encoded result as a byte slice — for example, storing in a database, sending over a message queue, or passing to another function.

### StreamCodec[T] — io.Reader/io.Writer-oriented

```go
// StreamCodec encodes T to an io.Writer and decodes T from an io.Reader.
type StreamCodec[T any] interface {
    Encode(w io.Writer, v T) error
    Decode(r io.Reader, v *T) error
}
```

Use `StreamCodec[T]` when working with streams — HTTP request/response bodies, files, network connections, or any `io.Reader`/`io.Writer`.

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

All codecs in this library are safe for concurrent use. Serialization codecs like `json.Codec[T]` are stateless zero-size structs; compression wrappers hold only immutable configuration. You can safely share a single codec instance across goroutines.
