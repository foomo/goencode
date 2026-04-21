# File Codec

The `file` package wraps any `Codec[T, []byte]` or `StreamCodec[T]` to add atomic file persistence. It writes to a temporary file first, then renames it into place — preventing partial writes if the process crashes mid-write.

::: warning
The file codec has a different method signature — it uses `path string` instead of `[]byte` or `io.Writer`.
:::

## Basic Usage

```go
import (
    "github.com/foomo/goencode/file"
    "github.com/foomo/goencode/json/v1"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

fc := file.NewCodec[Config](json.NewCodec[Config]()) // [!code highlight]

// Write config to file (atomic: temp file + rename)
err := fc.Encode("config.json", Config{Host: "localhost", Port: 8080}) // [!code highlight]

// Read config from file
var cfg Config
err = fc.Decode("config.json", &cfg) // [!code highlight]
```

## Atomic Writes

When `Encode` is called:

1. Data is serialized using the inner codec
2. A temporary file is created in the same directory as the target
3. Serialized bytes are written to the temp file
4. The temp file is renamed to the target path

The rename is atomic on the same filesystem, so readers always see either the old file or the complete new file — never a partial write.

## File Permissions

Use `WithPermissions` to set the file mode. The default is `0o644`.

```go
// Restrict to owner-only for sensitive data
fc := file.NewCodec[Secrets](
    json.NewCodec[Secrets](),
    file.WithPermissions(0o600), // [!code highlight]
)
```

## Stream Variant

`file.NewStreamCodec[T]` wraps a `StreamCodec[T]` struct instead. This streams the encoded data directly to the temp file without buffering the full payload in memory.

```go
fsc := file.NewStreamCodec[Config](json.NewStreamCodec[Config]()) // [!code highlight]

err := fsc.Encode("config.json", cfg)

var loaded Config
err = fsc.Decode("config.json", &loaded)
```

## Composed Example

Combine serialization, compression, and file persistence via `PipeCodec`:

```go
import (
    "github.com/foomo/goencode"
    "github.com/foomo/goencode/file"
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type State struct {
    Version int                    `json:"version"`
    Data    map[string]interface{} `json:"data"`
}

fc := file.NewCodec[State](
    goencode.PipeCodec(json.NewCodec[State](), gzip.NewCodec(gzip.WithLevel(gzip.BestSpeed))),
    file.WithPermissions(0o600),
)

// Writes JSON → gzip → temp file → rename
err := fc.Encode("/var/lib/myapp/state.json.gz", state)

// Reads file → gunzip → JSON
var loaded State
err = fc.Decode("/var/lib/myapp/state.json.gz", &loaded)
```
