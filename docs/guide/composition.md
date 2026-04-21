# Composing Codecs

Codecs in goencode compose via `PipeCodec` — a type-safe function that chains two codecs together. Compression codecs are standalone `Codec[[]byte, []byte]`, so you pipe a serialization codec into a compression codec to get a single composed codec.

## Basic Composition

`PipeCodec` chains two codecs where the output type of the first matches the input type of the second:

```go
import (
    "github.com/foomo/goencode"
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// JSON serialization + gzip compression
c := goencode.PipeCodec(json.NewCodec[User](), gzip.NewCodec()) // [!code highlight]

b, err := c.Encode(User{Name: "Alice", Age: 30})
// b contains gzip-compressed JSON

var u User
err = c.Decode(b, &u)
// u == User{Name: "Alice", Age: 30}
```

The flow is: `Encode` serializes with the first codec, then compresses with the second. `Decode` decompresses with the second, then deserializes with the first.

## Choosing a Format

::: code-group

```go [JSON + gzip]
c := goencode.PipeCodec(json.NewCodec[User](), gzip.NewCodec())
```

```go [XML + flate]
c := goencode.PipeCodec(xml.NewCodec[User](), flate.NewCodec())
```

```go [Gob + snappy]
c := goencode.PipeCodec(gob.NewCodec[User](), snappy.NewCodec())
```

```go [JSON + zstd]
c := goencode.PipeCodec(json.NewCodec[User](), zstd.NewCodec())
```

:::

## Compression Options

gzip, flate, zstd, and brotli accept options to tune compression level:

```go
// gzip with best compression
c := goencode.PipeCodec(
    json.NewCodec[User](),
    gzip.NewCodec(gzip.WithLevel(gzip.BestCompression)), // [!code highlight]
)

// flate with best speed
c := goencode.PipeCodec(
    json.NewCodec[User](),
    flate.NewCodec(flate.WithLevel(flate.BestSpeed)), // [!code highlight]
)

// zstd with best compression
c := goencode.PipeCodec(
    json.NewCodec[User](),
    zstd.NewCodec(zstd.WithLevel(zstd.SpeedBestCompression)), // [!code highlight]
)
```

## Chaining Multiple Codecs

`PipeCodec` returns a `Codec`, so you can chain more than two:

```go
// JSON → gzip → base64
c := goencode.PipeCodec(
    goencode.PipeCodec(json.NewCodec[User](), gzip.NewCodec()),
    base64.NewCodec(),
)
```

## Adding File Persistence

The `file` codec wraps any `Codec[T, []byte]` to read and write files atomically:

```go
import (
    "github.com/foomo/goencode"
    "github.com/foomo/goencode/file"
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

// JSON + gzip + atomic file I/O
fc := file.NewCodec[Config]( // [!code highlight]
    goencode.PipeCodec(json.NewCodec[Config](), gzip.NewCodec()), // [!code highlight]
    file.WithPermissions(0o600),
) // [!code highlight]

// Write atomically (temp file + rename)
err := fc.Encode("/etc/myapp/config.json.gz", cfg)

// Read back
var loaded Config
err = fc.Decode("/etc/myapp/config.json.gz", &loaded)
```

The composition layers: `file` wraps the piped codec (`json` → `gzip`).

## Standalone Encoder/Decoder Composition

You can also compose individual encoder and decoder functions:

```go
// Compose encoders: User → []byte → []byte
enc := goencode.PipeEncoder(json.Encoder[User], gzip.NewEncoder())

// Compose decoders: []byte → []byte → User
dec := goencode.PipeDecoder(json.Decoder[User], gzip.NewDecoder())
```
