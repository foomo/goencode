# Composing Codecs

Compression codecs in goencode follow the decorator pattern — they wrap an inner `Codec[T]` to add a compression layer. This lets you compose serialization and compression in a single line.

## Basic Composition

A compression codec takes any `Codec[T]` as its first argument:

```go
import (
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// JSON serialization + gzip compression
c := gzip.NewCodec[User](json.NewCodec[User]()) // [!code highlight]

b, err := c.Encode(User{Name: "Alice", Age: 30})
// b contains gzip-compressed JSON

var u User
err = c.Decode(b, &u)
// u == User{Name: "Alice", Age: 30}
```

The flow is: `Encode` serializes with the inner codec, then compresses. `Decode` decompresses, then deserializes.

## Choosing a Format

::: code-group

```go [JSON + gzip]
c := gzip.NewCodec[User](json.NewCodec[User]())
```

```go [XML + flate]
c := flate.NewCodec[User](xml.NewCodec[User]())
```

```go [Gob + snappy]
c := snappy.NewCodec[User](gob.NewCodec[User]())
```

```go [JSON + zstd]
c := zstd.NewCodec[User](json.NewCodec[User]())
```

:::

## Compression Options

gzip, flate, and zstd accept options to tune compression level:

```go
// gzip with best compression
c := gzip.NewCodec[User](
    json.NewCodec[User](),
    gzip.WithLevel(gzip.BestCompression), // [!code highlight]
)

// flate with best speed
c := flate.NewCodec[User](
    json.NewCodec[User](),
    flate.WithLevel(flate.BestSpeed), // [!code highlight]
)

// zstd with best compression
c := zstd.NewCodec[User](
    json.NewCodec[User](),
    zstd.WithLevel(zstd.SpeedBestCompression), // [!code highlight]
)
```

## Adding File Persistence

The `file` codec wraps any `Codec[T]` to read and write files atomically:

```go
import (
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
    gzip.NewCodec[Config](json.NewCodec[Config]()), // [!code highlight]
    file.WithPermissions(0o600),
) // [!code highlight]

// Write atomically (temp file + rename)
err := fc.Encode("/etc/myapp/config.json.gz", cfg)

// Read back
var loaded Config
err = fc.Decode("/etc/myapp/config.json.gz", &loaded)
```

The composition layers from outside in: `file` → `gzip` → `json` → `T`.

## Stream Composition

The same pattern works with `StreamCodec[T]`:

```go
import (
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

sc := gzip.NewStreamCodec[User](json.NewStreamCodec[User]()) // [!code highlight]

// Write compressed JSON to any io.Writer
err := sc.Encode(writer, user)

// Read compressed JSON from any io.Reader
var u User
err = sc.Decode(reader, &u)
```
