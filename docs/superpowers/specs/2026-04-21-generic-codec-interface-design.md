# Generic Codec Interface Redesign

**Date:** 2026-04-21
**Status:** Draft

## Summary

Redesign goencode's core interfaces from single-param `Codec[T]` / `StreamCodec[T]` interfaces to two-param `Codec[S, T]` / `StreamCodec[S]` function-type structs. Enables type-safe codec composition via `Pipe`, unifies encoding/compression/file/conversion codecs under one model.

## Motivation

- Current `Codec[T]` hardcodes `[]byte` as target — file codec (`string` path) and base64 (`[]byte → []byte`) are special cases that don't fit the interface
- No way to compose codecs with type-safe piping (e.g., JSON → base64 → file)
- Type conversion codecs (e.g., `string ↔ int`) impossible under current interface
- Compression codecs unnecessarily coupled to inner codec via decorator pattern

## Core Types

### Primitives

```go
package goencode

import "io"

// Encoder encodes source S to target T.
type Encoder[S, T any] func(s S) (T, error)

// Decoder decodes target T back into source S.
type Decoder[S, T any] func(t T, s *S) error

// StreamEncoder encodes S into an io.Writer.
type StreamEncoder[S any] func(w io.Writer, s S) error

// StreamDecoder decodes S from an io.Reader.
type StreamDecoder[S any] func(r io.Reader, s *S) error
```

### Bundles

```go
// Codec bundles an Encoder and Decoder for S ↔ T round-trips.
type Codec[S, T any] struct {
    Encode Encoder[S, T]
    Decode Decoder[S, T]
}

// StreamCodec bundles streaming encode/decode for S.
type StreamCodec[S any] struct {
    Encode StreamEncoder[S]
    Decode StreamDecoder[S]
}
```

### Design Decisions

- **Function types over interfaces**: most composable, least boilerplate. Closures capture config naturally.
- **Pointer baked into Decoder**: `Decoder[S, T]` signature is `func(t T, s *S) error` — pointer on S is implicit, avoids third type param.
- **StreamCodec stays single type param**: io.Writer/io.Reader are fixed, no need for `[S, T]`.
- **Structs not interfaces**: Codec/StreamCodec are struct bundles of function fields, not interface contracts.

## Composition

```go
// PipeEncoder chains two encoders: A → B → C.
func PipeEncoder[A, B, C any](first Encoder[A, B], second Encoder[B, C]) Encoder[A, C] {
    return func(a A) (C, error) {
        b, err := first(a)
        if err != nil {
            var zero C
            return zero, err
        }
        return second(b)
    }
}

// PipeDecoder chains two decoders: C → B → A (reverse order).
func PipeDecoder[A, B, C any](first Decoder[A, B], second Decoder[B, C]) Decoder[A, C] {
    return func(c C, a *A) error {
        var b B
        if err := second(c, &b); err != nil {
            return err
        }
        return first(b, a)
    }
}

// PipeCodec chains two codecs: Codec[A,B] + Codec[B,C] → Codec[A,C].
func PipeCodec[A, B, C any](first Codec[A, B], second Codec[B, C]) Codec[A, C] {
    return Codec[A, C]{
        Encode: PipeEncoder(first.Encode, second.Encode),
        Decode: PipeDecoder(first.Decode, second.Decode),
    }
}
```

## Codec Migration

### Serialization codecs (json, xml, gob, asn1, csv)

```go
// Constructor signature unchanged, return type changes
func NewCodec[T any]() goencode.Codec[T, []byte]
func NewStreamCodec[T any]() goencode.StreamCodec[T]
```

### Encoding codecs (base64, base32, hex, ascii85, pem)

```go
// Now fits naturally as []byte → []byte
func NewCodec() goencode.Codec[[]byte, []byte]
func NewStreamCodec() goencode.StreamCodec[[]byte]
```

### Compression codecs (gzip, flate, snappy, zstd)

No longer decorators. Become standalone `Codec[[]byte, []byte]`, compose via Pipe.

```go
// Before (decorator)
c := gzip.NewCodec(json.NewCodec[MyType]())

// After (Pipe composition)
c := goencode.PipeCodec(json.NewCodec[MyType](), gzip.NewCodec())
```

```go
func NewCodec(opts ...Option) goencode.Codec[[]byte, []byte]
func NewStreamCodec(opts ...Option) goencode.StreamCodec[[]byte]
```

### File codec

Becomes standalone `Codec[[]byte, string]` — compose via Pipe like compression codecs.

```go
// []byte ↔ string (file path). Encode writes bytes to file, returns path. Decode reads file.
func NewCodec(opts ...Option) goencode.Codec[[]byte, string]
```

```go
// Usage: JSON → file via Pipe
full := goencode.PipeCodec(json.NewCodec[MyType](), file.NewCodec()) // Codec[MyType, string]
```

Note: Encode requires caller to provide the file path. Signature is `func(b []byte) (string, error)` — but file codec needs a path to write to. Options: pass path via `WithPath(p string)` option, or change to `Codec[[]byte, string]` where encode takes bytes and an option sets the target path. Alternative: keep decorator pattern for file codec since it needs write path context. **Decision: keep file codec as decorator** — it wraps an inner codec because it needs to control the full write-path lifecycle (temp file + rename). Unlike compression, file I/O is inherently stateful (needs path).

```go
// File codec stays as wrapper — needs path context
func NewCodec[T any](codec goencode.Codec[T, []byte], opts ...Option) *Codec[T]

type Codec[T any] struct {
    codec goencode.Codec[T, []byte]
    perm  os.FileMode
}

// Encode writes v to file at path atomically. Decode reads file at path into v.
func (c *Codec[T]) Encode(path string, v T) error
func (c *Codec[T]) Decode(path string, v *T) error
```

File codec does NOT implement `goencode.Codec[S, T]` — it has its own signature with path. This is acceptable: file I/O is fundamentally different from value transformations.

### Type conversion codecs (new)

New capability enabled by `Codec[S, T]`:

```go
func NewStringIntCodec() goencode.Codec[string, int]
```

## Composition Examples

```go
// JSON → base64 encoded
jsonCodec := json.NewCodec[MyType]()                    // Codec[MyType, []byte]
b64Codec := base64.NewCodec()                           // Codec[[]byte, []byte]
combined := goencode.PipeCodec(jsonCodec, b64Codec)     // Codec[MyType, []byte]

// JSON → gzip compressed
jsonCodec := json.NewCodec[MyType]()                    // Codec[MyType, []byte]
gzipCodec := gzip.NewCodec()                            // Codec[[]byte, []byte]
combined := goencode.PipeCodec(jsonCodec, gzipCodec)    // Codec[MyType, []byte]

// JSON → base64, then write to file
jsonB64 := goencode.PipeCodec(json.NewCodec[MyType](), base64.NewCodec()) // Codec[MyType, []byte]
fc := file.NewCodec(jsonB64)                                               // file.Codec[MyType]
fc.Encode("/tmp/data.b64", myVal)                                          // atomic write

// Type conversion chaining
strToInt := conv.NewStringIntCodec()                    // Codec[string, int]
intToBytes := conv.NewIntBytesCodec()                   // Codec[int, []byte]
strToBytes := goencode.PipeCodec(strToInt, intToBytes)  // Codec[string, []byte]
```

## Removed Types

| Old | Replacement |
|-----|-------------|
| `Codec[T]` interface | `Codec[S, T]` struct |
| `StreamCodec[T]` interface | `StreamCodec[S]` struct |
| `Encoder[T]` interface | `Encoder[S, T]` func type |
| `Decoder[T]` interface | `Decoder[S, T]` func type |
| `EncoderFunc[T]` | Redundant — Encoder is already a func type |
| `DecoderFunc[T]` | Redundant — Decoder is already a func type |
| `StreamEncoder[T]` interface | `StreamEncoder[S]` func type |
| `StreamDecoder[T]` interface | `StreamDecoder[S]` func type |
| `StreamEncoderFunc[T]` | Redundant |
| `StreamDecoderFunc[T]` | Redundant |

## Breaking Changes

- All codec constructors return structs instead of interfaces
- Compression codecs no longer take inner codec — compose via `PipeCodec`
- File codec returns `Codec[T, string]` instead of custom interface
- All root package types replaced (see table above)
- Still alpha — no major version bump needed

## Testing Strategy

- **Round-trip tests** for every codec: encode then decode, assert equality
- **Pipe tests**: chain 2-3 codecs, verify round-trip through full pipeline
- **Error propagation**: first encoder fails → second never called
- **Decode reversal**: verify PipeDecoder applies decoders in reverse order
- **StreamCodec tests**: same pattern as today, unchanged
- **Benchmarks**: existing benchmark structure stays, update signatures
