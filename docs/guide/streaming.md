# Streaming

`StreamCodec[T]` encodes to an `io.Writer` and decodes from an `io.Reader`. Use it when working with streams — HTTP bodies, files, network connections, or pipelines — to avoid buffering entire payloads in memory.

## Interface

```go
type StreamCodec[T any] interface {
    Encode(w io.Writer, v T) error
    Decode(r io.Reader, v *T) error
}
```

The root package also defines `Encoder[T]` and `Decoder[T]` for stateful encoder/decoder pairs:

```go
type Encoder[T any] interface {
    Encode(v T) error
}

type Decoder[T any] interface {
    Decode(v any) error
}
```

## JSON Streaming

```go
import (
    "bytes"
    "fmt"

    "github.com/foomo/goencode/json/v1"
)

type Data struct {
    Name string
}

c := json.NewStreamCodec[Data]()

var buf bytes.Buffer
if err := c.Encode(&buf, Data{Name: "example-123"}); err != nil { // [!code highlight]
    log.Fatal(err)
}

var decoded Data
if err := c.Decode(&buf, &decoded); err != nil { // [!code highlight]
    log.Fatal(err)
}
fmt.Println(decoded.Name) // example-123
```

::: tip
`bytes.Buffer` implements both `io.Writer` and `io.Reader`, making it convenient for testing and in-memory pipelines.
:::

## Base64 Streaming

```go
import (
    "bytes"
    "fmt"

    "github.com/foomo/goencode/base64"
)

c := base64.NewStreamCodec()

var buf bytes.Buffer
_ = c.Encode(&buf, []byte("hello"))
fmt.Println(buf.String()) // aGVsbG8=

var decoded []byte
_ = c.Decode(&buf, &decoded)
fmt.Println(string(decoded)) // hello
```

## CSV Streaming

```go
import (
    "bytes"
    "fmt"

    "github.com/foomo/goencode/csv"
)

c := csv.NewStreamCodec()

records := [][]string{
    {"name", "age"},
    {"Alice", "30"},
}

var buf bytes.Buffer
_ = c.Encode(&buf, records) // [!code highlight]

var decoded [][]string
_ = c.Decode(&buf, &decoded) // [!code highlight]
fmt.Println(decoded) // [[name age] [Alice 30]]
```

## Compressed Streams

Stream codecs compose the same way as byte codecs — compression wrappers accept an inner `StreamCodec[T]`:

```go
import (
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type Payload struct {
    Items []string `json:"items"`
}

sc := gzip.NewStreamCodec[Payload](json.NewStreamCodec[Payload]()) // [!code highlight]

// Write gzip-compressed JSON to any io.Writer
err := sc.Encode(writer, payload)

// Read gzip-compressed JSON from any io.Reader
var p Payload
err = sc.Decode(reader, &p)
```

## HTTP Example

StreamCodec is a natural fit for HTTP handlers:

```go
func handleEncode(w http.ResponseWriter, r *http.Request) {
    sc := json.NewStreamCodec[Response]()
    w.Header().Set("Content-Type", "application/json")
    _ = sc.Encode(w, Response{Status: "ok"}) // [!code highlight]
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
    sc := json.NewStreamCodec[Request]()
    var req Request
    if err := sc.Decode(r.Body, &req); err != nil { // [!code highlight]
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // use req...
}
```
