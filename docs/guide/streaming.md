# Streaming

`StreamCodec[S]` encodes to an `io.Writer` and decodes from an `io.Reader`. Use it when working with streams — HTTP bodies, files, network connections, or pipelines — to avoid buffering entire payloads in memory.

## Types

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

Each serialization package exports standalone `StreamEncoder` and `StreamDecoder` functions alongside the `NewStreamCodec` constructor.

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

## Compression Streaming

Compression stream codecs are standalone `StreamCodec[[]byte]` — they compress and decompress raw bytes over streams:

```go
import (
    "github.com/foomo/goencode/gzip"
)

sc := gzip.NewStreamCodec() // StreamCodec[[]byte] // [!code highlight]

// Write gzip-compressed bytes to any io.Writer
err := sc.Encode(writer, rawBytes)

// Read gzip-decompressed bytes from any io.Reader
var decoded []byte
err = sc.Decode(reader, &decoded)
```

::: tip
For combined serialization + compression (e.g., JSON → gzip), use the byte-oriented `PipeCodec` approach instead. See [Composing Codecs](/guide/composition) for details.
:::

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
