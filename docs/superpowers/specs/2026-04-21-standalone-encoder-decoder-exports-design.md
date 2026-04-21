# Standalone Encoder/Decoder Exports

**Date:** 2026-04-21
**Status:** Approved

## Problem

Consumer APIs (e.g. a messaging library) often need only one direction — decode incoming messages or encode outgoing ones. Currently most subpackages only export `NewCodec[T]()` returning a full `Codec[S, T]`, forcing consumers to depend on both directions even when they only need one.

The root package already defines `Encoder[S, T]` and `Decoder[S, T]` as standalone function types, and `json/v1` already exports bare `Encoder[T]`/`Decoder[T]` funcs. This pattern should be extended to all subpackages.

## Design

### Rule

- **No options** → export bare funcs `Encoder` and `Decoder` (like `json/v1` today)
- **Takes options** → export `NewEncoder(opts ...Option)` and `NewDecoder(opts ...Option)` constructors

### Simple codecs — bare funcs

Serialization codecs (generic `[T any]`):

| Package | Exports |
|---------|---------|
| json/v1 | `Encoder[T]`, `Decoder[T]` *(already exists)* |
| json/v2 | `Encoder[T]`, `Decoder[T]` |
| xml | `Encoder[T]`, `Decoder[T]` |
| gob | `Encoder[T]`, `Decoder[T]` |
| asn1 | `Encoder[T]`, `Decoder[T]` |
| csv | `Encoder[T]`, `Decoder[T]` |
| toml | `Encoder[T]`, `Decoder[T]` |
| yaml/v2 | `Encoder[T]`, `Decoder[T]` |
| yaml/v3 | `Encoder[T]`, `Decoder[T]` |
| yaml/v4 | `Encoder[T]`, `Decoder[T]` |
| msgpack/tinylib | `Encoder[T]`, `Decoder[T]` |
| msgpack/vmihailenco | `Encoder[T]`, `Decoder[T]` |

Encoding codecs (no type param, `[]byte` ↔ `[]byte` or `*pem.Block` ↔ `[]byte`):

| Package | Exports |
|---------|---------|
| base64 | `Encoder`, `Decoder` |
| base32 | `Encoder`, `Decoder` |
| hex | `Encoder`, `Decoder` |
| ascii85 | `Encoder`, `Decoder` |
| pem | `Encoder`, `Decoder` |
| snappy | `Encoder`, `Decoder` |

### Configurable codecs — constructor funcs

Compression codecs that accept options:

| Package | Exports |
|---------|---------|
| gzip | `NewEncoder(opts ...Option)`, `NewDecoder(opts ...Option)` |
| flate | `NewEncoder(opts ...Option)`, `NewDecoder(opts ...Option)` |
| zstd | `NewEncoder(opts ...Option)`, `NewDecoder(opts ...Option)` |
| brotli | `NewEncoder(opts ...Option)`, `NewDecoder(opts ...Option)` |

### Skipped

- `file` — wraps another codec + filesystem I/O. Different abstraction, leave as-is.

## Call-site examples

Consumer API requiring only decode:

```go
func NewConsumer[T any](decode goencode.Decoder[T, []byte]) { ... }

NewConsumer(json.Decoder[MyMsg])
NewConsumer(gob.Decoder[MyMsg])
NewConsumer(yaml.Decoder[MyMsg])
```

Compression — only encode direction:

```go
compress := gzip.NewEncoder(gzip.WithLevel(gzip.BestCompression))
data, err := compress(raw)
```

## Scope

- Purely additive — existing `NewCodec`/`NewStreamCodec` constructors unchanged.
- Each subpackage's `NewCodec` reuses the new bare funcs / constructors internally.
- No changes to root package types (`Encoder`, `Decoder`, `Codec`, `Pipe*`).
- StreamEncoder/StreamDecoder standalone exports are out of scope for this change.
