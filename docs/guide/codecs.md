# Codec Reference

All codecs listed below implement `Codec[T]`, `StreamCodec[T]`, or both. The core module packages use only the Go standard library.

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

## Compression Wrappers

Compression codecs wrap an inner `Codec[T]` or `StreamCodec[T]` using the [decorator pattern](/guide/composition).

| Package | Constructor | Options |
|---------|-------------|---------|
| `gzip` | `gzip.NewCodec[T](codec, opts...)` | `gzip.WithLevel(int)` |
| `flate` | `flate.NewCodec[T](codec, opts...)` | `flate.WithLevel(int)` |
| `snappy` | `snappy.NewCodec[T](codec)` | — |
| `zstd` | `zstd.NewCodec[T](codec, opts...)` | `zstd.WithLevel(zstd.EncoderLevel)` |
| `brotli` | `brotli.NewCodec[T](codec, opts...)` | `brotli.WithLevel(int)` |

Each also has a stream variant: `gzip.NewStreamCodec[T](streamCodec, opts...)`, etc.

::: tip
`snappy`, `zstd`, and `brotli` are [submodule packages](#submodule-packages) that require a separate `go get`.
:::


## Utility

### File Codec

The `file` package wraps any codec to read/write files atomically (temp file + rename).

```go
file.NewCodec[T](codec, opts...)       // Encode(path string, v T) error
file.NewStreamCodec[T](codec, opts...) // Encode(path string, v T) error
```

Options: `file.WithPermissions(os.FileMode)` — default `0o644`.

::: warning
The file codec has a different method signature — it uses `path string` instead of `[]byte` or `io.Writer`. It does not satisfy the `Codec[T]` or `StreamCodec[T]` interfaces. See [File Codec](/guide/file-codec) for details.
:::

## Submodule Packages

These packages have external dependencies and live in separate Go modules. Install them individually.

| Package | Import Path | Dependency | StreamCodec |
|---------|-------------|------------|-------------|
| `json2` | `github.com/foomo/goencode/json2` | go-json-experiment | — (has `EncodeTo`/`DecodeFrom` methods) |
| `yaml/v2` | `github.com/foomo/goencode/yaml/v2` | go.yaml.in/yaml/v2 | — |
| `yaml/v3` | `github.com/foomo/goencode/yaml/v3` | go.yaml.in/yaml/v3 | — |
| `yaml/v4` | `github.com/foomo/goencode/yaml/v4` | go.yaml.in/yaml/v4 | — |
| `snappy` | `github.com/foomo/goencode/snappy` | github.com/golang/snappy | `snappy.NewStreamCodec[T](codec)` |
| `brotli` | `github.com/foomo/goencode/brotli` | github.com/andybalholm/brotli | `brotli.NewStreamCodec[T](codec, opts...)` |
| `toml` | `github.com/foomo/goencode/toml` | github.com/BurntSushi/toml | `toml.NewStreamCodec[T]()` |
| `zstd` | `github.com/foomo/goencode/zstd` | github.com/klauspost/compress | `zstd.NewStreamCodec[T](codec, opts...)` |
