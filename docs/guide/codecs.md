# Codec Reference

All codecs listed below return `Codec[S, T]` and/or `StreamCodec[S]` structs via constructor functions. Each package also exports standalone `Encoder` and `Decoder` functions. The core module packages use only the Go standard library.

## Serialization Codecs

| Package | Constructor | Type | StreamCodec |
|---------|-------------|------|-------------|
| `json` | `json.NewCodec[T]()` | generic `T` | `json.NewStreamCodec[T]()` |
| `xml` | `xml.NewCodec[T]()` | generic `T` | `xml.NewStreamCodec[T]()` |
| `gob` | `gob.NewCodec[T]()` | generic `T` | `gob.NewStreamCodec[T]()` |
| `asn1` | `asn1.NewCodec[T]()` | generic `T` | `asn1.NewStreamCodec[T]()` |
| `csv` | `csv.NewCodec()` | `[][]string` | `csv.NewStreamCodec()` |
| `pem` | `pem.NewCodec()` | `*pem.Block` | `pem.NewStreamCodec()` |

## Binary Encoding Codecs

| Package | Constructor | Type | StreamCodec |
|---------|-------------|------|-------------|
| `base64` | `base64.NewCodec()` | `[]byte` | `base64.NewStreamCodec()` |
| `base32` | `base32.NewCodec()` | `[]byte` | `base32.NewStreamCodec()` |
| `hex` | `hex.NewCodec()` | `[]byte` | `hex.NewStreamCodec()` |
| `ascii85` | `ascii85.NewCodec()` | `[]byte` | `ascii85.NewStreamCodec()` |

## Compression Codecs

Compression codecs are standalone `Codec[[]byte, []byte]` — they compress and decompress raw bytes. Compose them with serialization codecs via [`PipeCodec`](/guide/composition).

| Package | Constructor | StreamCodec | Options |
|---------|-------------|-------------|---------|
| `gzip` | `gzip.NewCodec(opts...)` | `gzip.NewStreamCodec(opts...)` | `gzip.WithLevel(int)` |
| `flate` | `flate.NewCodec(opts...)` | `flate.NewStreamCodec(opts...)` | `flate.WithLevel(int)` |
| `snappy` | `snappy.NewCodec()` | `snappy.NewStreamCodec()` | — |
| `zstd` | `zstd.NewCodec(opts...)` | `zstd.NewStreamCodec(opts...)` | `zstd.WithLevel(zstd.EncoderLevel)` |
| `brotli` | `brotli.NewCodec(opts...)` | `brotli.NewStreamCodec(opts...)` | `brotli.WithLevel(int)` |

```go
// Example: JSON + gzip via PipeCodec
c := goencode.PipeCodec(json.NewCodec[User](), gzip.NewCodec())
```

::: tip
`snappy`, `zstd`, and `brotli` are [submodule packages](#submodule-packages) that require a separate `go get`.
:::


## Utility

### File Codec

The `file` package wraps any codec to read/write files atomically (temp file + rename).

```go
file.NewCodec[T](codec, opts...)       // accepts Codec[T, []byte], Encode(path string, v T) error
file.NewStreamCodec[T](codec, opts...) // accepts StreamCodec[T], Encode(path string, v T) error
```

Options: `file.WithPermissions(os.FileMode)` — default `0o644`.

::: warning
The file codec has a different method signature — it uses `path string` instead of `[]byte` or `io.Writer`. See [File Codec](/guide/file-codec) for details.
:::

## Submodule Packages

These packages have external dependencies and live in separate Go modules. Install them individually.

| Package | Import Path | Dependency | StreamCodec |
|---------|-------------|------------|-------------|
| `json/v2` | `github.com/foomo/goencode/json/v2` | go-json-experiment | `json.NewStreamCodec[T]()` |
| `yaml/v2` | `github.com/foomo/goencode/yaml/v2` | go.yaml.in/yaml/v2 | `yaml.NewStreamCodec[T]()` |
| `yaml/v3` | `github.com/foomo/goencode/yaml/v3` | go.yaml.in/yaml/v3 | `yaml.NewStreamCodec[T]()` |
| `yaml/v4` | `github.com/foomo/goencode/yaml/v4` | go.yaml.in/yaml/v4 | `yaml.NewStreamCodec[T]()` |
| `toml` | `github.com/foomo/goencode/toml` | github.com/BurntSushi/toml | `toml.NewStreamCodec[T]()` |
| `msgpack/tinylib` | `github.com/foomo/goencode/msgpack/tinylib` | github.com/tinylib/msgp | `msgpack.NewStreamCodec[T]()` |
| `msgpack/vmihailenco` | `github.com/foomo/goencode/msgpack/vmihailenco` | github.com/vmihailenco/msgpack | `msgpack.NewStreamCodec[T]()` |
| `snappy` | `github.com/foomo/goencode/snappy` | github.com/golang/snappy | `snappy.NewStreamCodec()` |
| `brotli` | `github.com/foomo/goencode/brotli` | github.com/andybalholm/brotli | `brotli.NewStreamCodec(opts...)` |
| `zstd` | `github.com/foomo/goencode/zstd` | github.com/klauspost/compress | `zstd.NewStreamCodec(opts...)` |
