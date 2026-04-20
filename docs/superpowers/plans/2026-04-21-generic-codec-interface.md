# Generic Codec Interface Redesign — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace single-param `Codec[T]`/`StreamCodec[T]` interfaces with two-param `Codec[S, T]`/`StreamCodec[S]` function-type structs, enabling type-safe codec composition via `Pipe`.

**Architecture:** Core types become function types (`Encoder[S, T]`, `Decoder[S, T]`) bundled into structs (`Codec[S, T]`, `StreamCodec[S]`). Composition via free `Pipe*` functions. Compression codecs become standalone `Codec[[]byte, []byte]` instead of decorators. File codec stays as decorator with own signature.

**Tech Stack:** Go 1.24+, generics, no new dependencies.

**Spec:** `docs/superpowers/specs/2026-04-21-generic-codec-interface-design.md`

---

## File Structure

### Root package (`github.com/foomo/goencode`)

| Action | File | Responsibility |
|--------|------|---------------|
| Rewrite | `codec.go` | `Codec[S, T]` struct, `Encoder[S, T]` func type, `Decoder[S, T]` func type |
| Rewrite | `streamcodec.go` | `StreamCodec[S]` struct, `StreamEncoder[S]` func type, `StreamDecoder[S]` func type |
| Create | `pipe.go` | `PipeEncoder`, `PipeDecoder`, `PipeCodec` functions |
| Delete | `encoder.go` | Old `Encoder[T]` interface — replaced by `Encoder[S, T]` func type in codec.go |
| Delete | `decoder.go` | Old `Decoder[T]` interface — replaced by `Decoder[S, T]` func type in codec.go |
| Delete | `encoderfunc.go` | Old `EncoderFunc[T]` — redundant |
| Delete | `decoderfunc.go` | Old `DecoderFunc[T]` — redundant |
| Delete | `streamencoder.go` | Old `StreamEncoder[T]` interface — replaced |
| Delete | `streamdecoder.go` | Old `StreamDecoder[T]` interface — replaced |
| Delete | `streamencoderfunc.go` | Old `StreamEncoderFunc[T]` — redundant |
| Delete | `streamdecoderfunc.go` | Old `StreamDecoderFunc[T]` — redundant |

### Subpackages (each follows same pattern)

| Category | Packages | Codec change | StreamCodec change |
|----------|----------|-------------|-------------------|
| Serialization | `json/v1`, `json/v2`, `xml`, `gob`, `asn1` | Return `goencode.Codec[T, []byte]` | Return `goencode.StreamCodec[T]` |
| Encoding | `base64`, `base32`, `hex`, `ascii85`, `pem` | Return `goencode.Codec[[]byte, []byte]` | Return `goencode.StreamCodec[[]byte]` |
| CSV | `csv` | Return `goencode.Codec[[][]string, []byte]` | Return `goencode.StreamCodec[[][]string]` |
| Compression | `gzip`, `flate`, `snappy`, `zstd`, `brotli` | Standalone `goencode.Codec[[]byte, []byte]` (no inner codec) | Standalone `goencode.StreamCodec[[]byte]` (no inner codec) |
| File | `file` | Keeps own `Codec[T]` struct wrapping `goencode.Codec[T, []byte]` | Keeps own `StreamCodec[T]` wrapping `goencode.StreamCodec[T]` |
| YAML | `yaml/v2`, `yaml/v3`, `yaml/v4` | Return `goencode.Codec[T, []byte]` | N/A (no stream codecs) |
| Msgpack | `msgpack/tinylib`, `msgpack/vmihailenco` | Return `goencode.Codec[T, []byte]` | Return `goencode.StreamCodec[T]` |

---

### Task 1: Rewrite Root Package Core Types

**Files:**
- Rewrite: `codec.go`
- Rewrite: `streamcodec.go`
- Delete: `encoder.go`, `decoder.go`, `encoderfunc.go`, `decoderfunc.go`, `streamencoder.go`, `streamdecoder.go`, `streamencoderfunc.go`, `streamdecoderfunc.go`

- [ ] **Step 1: Delete obsolete files**

```bash
cd /Users/franklin/Workingcopies/github.com/foomo/goencode
rm encoder.go decoder.go encoderfunc.go decoderfunc.go \
   streamencoder.go streamdecoder.go streamencoderfunc.go streamdecoderfunc.go
```

- [ ] **Step 2: Rewrite `codec.go`**

```go
package goencode

// Encoder encodes source S to target T.
type Encoder[S, T any] func(s S) (T, error)

// Decoder decodes target T back into source S.
type Decoder[S, T any] func(t T, s *S) error

// Codec bundles an Encoder and Decoder for S ↔ T round-trips.
type Codec[S, T any] struct {
	Encode Encoder[S, T]
	Decode Decoder[S, T]
}
```

- [ ] **Step 3: Rewrite `streamcodec.go`**

```go
package goencode

import "io"

// StreamEncoder encodes S into an io.Writer.
type StreamEncoder[S any] func(w io.Writer, s S) error

// StreamDecoder decodes S from an io.Reader.
type StreamDecoder[S any] func(r io.Reader, s *S) error

// StreamCodec bundles streaming encode/decode for S.
type StreamCodec[S any] struct {
	Encode StreamEncoder[S]
	Decode StreamDecoder[S]
}
```

- [ ] **Step 4: Verify root package compiles**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go build ./...`

Expected: Compilation errors in subpackages (they still reference old types). Root package itself should compile.

Note: Use `go build .` (root only) to verify just the root package, since subpackages will fail until migrated.

Run: `go build .`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: replace Codec[T] interface with Codec[S,T] function-type struct

Replace single-param interfaces with two-param function types:
- Encoder[S,T] func type, Decoder[S,T] func type
- Codec[S,T] struct bundling Encode/Decode
- StreamEncoder[S], StreamDecoder[S] func types
- StreamCodec[S] struct bundling stream Encode/Decode
- Remove old Encoder[T], Decoder[T] interfaces and Func wrappers"
```

---

### Task 2: Add Pipe Composition Functions

**Files:**
- Create: `pipe.go`
- Create: `pipe_test.go`

- [ ] **Step 1: Write `pipe_test.go`**

```go
package goencode_test

import (
	"fmt"
	"strconv"
	"testing"

	goencode "github.com/foomo/goencode"
)

func TestPipeEncoder(t *testing.T) {
	intToStr := goencode.Encoder[int, string](func(i int) (string, error) {
		return strconv.Itoa(i), nil
	})
	strToBytes := goencode.Encoder[string, []byte](func(s string) ([]byte, error) {
		return []byte(s), nil
	})

	piped := goencode.PipeEncoder(intToStr, strToBytes)

	got, err := piped(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "42" {
		t.Fatalf("got %q, want %q", string(got), "42")
	}
}

func TestPipeEncoder_FirstError(t *testing.T) {
	failing := goencode.Encoder[int, string](func(i int) (string, error) {
		return "", fmt.Errorf("encode failed")
	})
	second := goencode.Encoder[string, []byte](func(s string) ([]byte, error) {
		t.Fatal("second encoder should not be called")
		return nil, nil
	})

	piped := goencode.PipeEncoder(failing, second)

	_, err := piped(42)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPipeDecoder(t *testing.T) {
	strToInt := goencode.Decoder[int, string](func(s string, i *int) error {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = v
		return nil
	})
	bytesToStr := goencode.Decoder[string, []byte](func(b []byte, s *string) error {
		*s = string(b)
		return nil
	})

	piped := goencode.PipeDecoder(strToInt, bytesToStr)

	var got int
	if err := piped([]byte("42"), &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

func TestPipeCodec(t *testing.T) {
	intStr := goencode.Codec[int, string]{
		Encode: func(i int) (string, error) {
			return strconv.Itoa(i), nil
		},
		Decode: func(s string, i *int) error {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			*i = v
			return nil
		},
	}
	strBytes := goencode.Codec[string, []byte]{
		Encode: func(s string) ([]byte, error) {
			return []byte(s), nil
		},
		Decode: func(b []byte, s *string) error {
			*s = string(b)
			return nil
		},
	}

	piped := goencode.PipeCodec(intStr, strBytes)

	encoded, err := piped.Encode(42)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	if string(encoded) != "42" {
		t.Fatalf("encoded: got %q, want %q", string(encoded), "42")
	}

	var decoded int
	if err := piped.Decode(encoded, &decoded); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded != 42 {
		t.Fatalf("decoded: got %d, want 42", decoded)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe -run TestPipe -v .`
Expected: FAIL — `PipeEncoder`, `PipeDecoder`, `PipeCodec` not defined

- [ ] **Step 3: Write `pipe.go`**

```go
package goencode

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

// PipeDecoder chains two decoders in reverse: decodes C → B via second, then B → A via first.
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

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe -run TestPipe -v .`
Expected: PASS (all 4 tests)

- [ ] **Step 5: Commit**

```bash
git add pipe.go pipe_test.go && git commit -m "feat: add Pipe composition functions for Encoder, Decoder, Codec"
```

---

### Task 3: Migrate Serialization Codecs (json/v1, xml, gob, asn1)

These all follow the same pattern: stateless generic codec returning `Codec[T, []byte]`.

**Files:**
- Modify: `json/v1/codec.go`
- Modify: `json/v1/streamcodec.go`
- Modify: `json/v1/codec_test.go`
- Modify: `json/v1/streamcodec_test.go`
- Modify: `xml/codec.go` (same pattern)
- Modify: `xml/streamcodec.go` (same pattern, if exists)
- Modify: `gob/streamcodec.go`
- Modify: `asn1/codec.go`
- Modify: `asn1/streamcodec.go`
- Modify: all corresponding `*_test.go` and `benchmark_test.go` files

- [ ] **Step 1: Rewrite `json/v1/codec.go`**

```go
package json

import (
	"encoding/json"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns a JSON codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: func(v T) ([]byte, error) {
			return json.Marshal(v)
		},
		Decode: func(b []byte, v *T) error {
			return json.Unmarshal(b, v)
		},
	}
}
```

- [ ] **Step 2: Rewrite `json/v1/streamcodec.go`**

```go
package json

import (
	"encoding/json"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a JSON stream codec for T.
// It is safe for concurrent use.
func NewStreamCodec[T any]() encoding.StreamCodec[T] {
	return encoding.StreamCodec[T]{
		Encode: func(w io.Writer, v T) error {
			return json.NewEncoder(w).Encode(v)
		},
		Decode: func(r io.Reader, v *T) error {
			return json.NewDecoder(r).Decode(v)
		},
	}
}
```

- [ ] **Step 3: Update `json/v1/codec_test.go`**

```go
package json_test

import (
	"fmt"

	"github.com/foomo/goencode/json/v1"
)

func ExampleNewCodec() {
	type Data struct {
		Name string
	}

	c := json.NewCodec[Data]()

	encoded, err := c.Encode(Data{Name: "example-123"})
	if err != nil {
		fmt.Printf("Encode failed: %v\n", err)
		return
	}

	fmt.Printf("Encoded: %s\n", string(encoded))

	var decoded Data
	if err := c.Decode(encoded, &decoded); err != nil {
		fmt.Printf("Decode failed: %v\n", err)
		return
	}

	fmt.Printf("Decoded Name: %s\n", decoded.Name)
	// Output:
	// Encoded: {"Name":"example-123"}
	// Decoded Name: example-123
}
```

Note: Example function name changes from `ExampleCodec` to `ExampleNewCodec` because there is no longer a `Codec` exported type — only `NewCodec` constructor.

- [ ] **Step 4: Migrate xml, gob, asn1 codecs using same pattern**

For each package, replace the struct type + methods with a constructor returning `encoding.Codec[T, []byte]` or `encoding.StreamCodec[T]`.

**`xml/codec.go`** — same as json/v1 but uses `encoding/xml.Marshal`/`Unmarshal`. Note: current xml codec uses `bufpool.sync` for encoding — keep that optimization by using a closure that captures the pool usage.

**`gob/streamcodec.go`** — returns `encoding.StreamCodec[T]` using `gob.NewEncoder`/`gob.NewDecoder`.

**`asn1/codec.go`** — uses `encoding/asn1.Marshal`/`Unmarshal`. Note: asn1.Unmarshal returns `(rest []byte, err error)` — keep the existing discard-rest pattern.

Delete the old struct types (`Codec[T]`, `StreamCodec[T]`) from each package — they are replaced by the constructor functions.

- [ ] **Step 5: Run tests for all migrated packages**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./json/v1/... ./xml/... ./gob/... ./asn1/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add json/v1/ xml/ gob/ asn1/ && git commit -m "refactor: migrate json/v1, xml, gob, asn1 to Codec[S,T] function types"
```

---

### Task 4: Migrate Encoding Codecs (base64, base32, hex, ascii85, pem)

These are non-generic, operating on `[]byte` → `[]byte` (or `*pem.Block` for pem).

**Files:**
- Modify: `base64/codec.go`, `base64/streamcodec.go`
- Modify: `base32/codec.go`, `base32/streamcodec.go`
- Modify: `hex/codec.go`, `hex/streamcodec.go`
- Modify: `ascii85/codec.go`, `ascii85/streamcodec.go`
- Modify: `pem/streamcodec.go`
- Modify: all corresponding test and benchmark files

- [ ] **Step 1: Rewrite `base64/codec.go`**

```go
package base64

import (
	stdbase64 "encoding/base64"

	encoding "github.com/foomo/goencode"
)

// NewCodec returns a Base64 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: func(v []byte) ([]byte, error) {
			dst := make([]byte, stdbase64.StdEncoding.EncodedLen(len(v)))
			stdbase64.StdEncoding.Encode(dst, v)
			return dst, nil
		},
		Decode: func(b []byte, v *[]byte) error {
			dst := make([]byte, stdbase64.StdEncoding.DecodedLen(len(b)))
			n, err := stdbase64.StdEncoding.Decode(dst, b)
			if err != nil {
				return err
			}
			*v = dst[:n]
			return nil
		},
	}
}
```

- [ ] **Step 2: Rewrite `base64/streamcodec.go`**

```go
package base64

import (
	stdbase64 "encoding/base64"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a Base64 stream codec.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[]byte] {
	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, v []byte) error {
			enc := stdbase64.NewEncoder(stdbase64.StdEncoding, w)
			if _, err := enc.Write(v); err != nil {
				_ = enc.Close()
				return err
			}
			return enc.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			data, err := io.ReadAll(stdbase64.NewDecoder(stdbase64.StdEncoding, r))
			if err != nil {
				return err
			}
			*v = data
			return nil
		},
	}
}
```

- [ ] **Step 3: Migrate base32, hex, ascii85, pem using same pattern**

Each follows the same structure — replace struct + methods with constructor returning `encoding.Codec[[]byte, []byte]` and `encoding.StreamCodec[[]byte]`.

**pem** is special: operates on `*pem.Block` not `[]byte`. Returns `encoding.Codec[*pem.Block, []byte]` and `encoding.StreamCodec[*pem.Block]`.

Delete old struct types from each package.

- [ ] **Step 4: Update test files**

Rename example functions from `ExampleCodec`/`ExampleStreamCodec` to `ExampleNewCodec`/`ExampleNewStreamCodec` since the exported type is now the constructor, not the struct.

- [ ] **Step 5: Run tests**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./base64/... ./base32/... ./hex/... ./ascii85/... ./pem/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add base64/ base32/ hex/ ascii85/ pem/ && git commit -m "refactor: migrate encoding codecs (base64, base32, hex, ascii85, pem) to Codec[S,T]"
```

---

### Task 5: Migrate CSV Codec

CSV is special — operates on `[][]string`.

**Files:**
- Modify: `csv/streamcodec.go`
- Modify: `csv/codec_test.go`, `csv/streamcodec_test.go`, `csv/benchmark_test.go`

- [ ] **Step 1: Rewrite `csv/streamcodec.go`**

CSV only has a StreamCodec. Return `encoding.StreamCodec[[][]string]`.

```go
package csv

import (
	"encoding/csv"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a CSV stream codec for [][]string.
// It is safe for concurrent use.
func NewStreamCodec() encoding.StreamCodec[[][]string] {
	return encoding.StreamCodec[[][]string]{
		Encode: func(w io.Writer, v [][]string) error {
			return csv.NewWriter(w).WriteAll(v)
		},
		Decode: func(r io.Reader, v *[][]string) error {
			records, err := csv.NewReader(r).ReadAll()
			if err != nil {
				return err
			}
			*v = records
			return nil
		},
	}
}
```

Note: Check if csv also has a `Codec` (non-stream). If so, migrate it similarly to return `encoding.Codec[[][]string, []byte]`.

- [ ] **Step 2: Update test files, run tests**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./csv/...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add csv/ && git commit -m "refactor: migrate csv codec to StreamCodec[S] function type"
```

---

### Task 6: Migrate Compression Codecs (gzip, flate, snappy, zstd, brotli)

Biggest change: remove decorator pattern. Each becomes standalone `Codec[[]byte, []byte]`.

**Files:**
- Rewrite: `gzip/codec.go`, `gzip/streamcodec.go`
- Rewrite: `flate/codec.go`, `flate/streamcodec.go`
- Rewrite: `snappy/codec.go`, `snappy/streamcodec.go`
- Rewrite: `zstd/codec.go`, `zstd/streamcodec.go`
- Rewrite: `brotli/codec.go`, `brotli/streamcodec.go`
- Keep: `gzip/option.go`, `flate/option.go`, `zstd/option.go`, `brotli/option.go` (unchanged)
- Modify: all corresponding test and benchmark files

- [ ] **Step 1: Rewrite `gzip/codec.go`**

```go
package gzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewCodec returns a gzip compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	o := options{
		level: gzip.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.Codec[[]byte, []byte]{
		Encode: func(data []byte) ([]byte, error) {
			buf := sync.Get()
			defer sync.Put(buf)

			w, err := gzip.NewWriterLevel(buf, o.level)
			if err != nil {
				return nil, err
			}

			if _, err := w.Write(data); err != nil {
				return nil, err
			}

			if err := w.Close(); err != nil {
				return nil, err
			}

			return append([]byte(nil), buf.Bytes()...), nil
		},
		Decode: func(data []byte, v *[]byte) error {
			r, err := gzip.NewReader(bytes.NewReader(data))
			if err != nil {
				return err
			}
			defer r.Close()

			var src io.Reader = r
			if o.maxDecodedSize > 0 {
				src = io.LimitReader(r, o.maxDecodedSize+1)
			}

			decoded, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			if o.maxDecodedSize > 0 && int64(len(decoded)) > o.maxDecodedSize {
				return fmt.Errorf("gzip: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
			}

			*v = decoded
			return nil
		},
	}
}
```

- [ ] **Step 2: Rewrite `gzip/streamcodec.go`**

```go
package gzip

import (
	"compress/gzip"
	"io"

	encoding "github.com/foomo/goencode"
)

// NewStreamCodec returns a gzip compression stream codec.
// It is safe for concurrent use.
func NewStreamCodec(opts ...Option) encoding.StreamCodec[[]byte] {
	o := options{
		level: gzip.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return encoding.StreamCodec[[]byte]{
		Encode: func(w io.Writer, data []byte) error {
			gw, err := gzip.NewWriterLevel(w, o.level)
			if err != nil {
				return err
			}

			if _, err := gw.Write(data); err != nil {
				gw.Close()
				return err
			}

			return gw.Close()
		},
		Decode: func(r io.Reader, v *[]byte) error {
			gr, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer gr.Close()

			var src io.Reader = gr
			if o.maxDecodedSize > 0 {
				src = io.LimitReader(gr, o.maxDecodedSize+1)
			}

			data, err := io.ReadAll(src)
			if err != nil {
				return err
			}

			if o.maxDecodedSize > 0 && int64(len(data)) > o.maxDecodedSize {
				return fmt.Errorf("gzip: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
			}

			*v = data
			return nil
		},
	}
}
```

- [ ] **Step 3: Migrate flate, snappy, zstd, brotli using same pattern**

Each compression codec follows the same transformation:
- Remove generic type param `[T any]`
- Remove inner `codec` field
- Return `encoding.Codec[[]byte, []byte]` / `encoding.StreamCodec[[]byte]`
- Encode/Decode operate directly on `[]byte`
- Keep option.go files unchanged

**snappy** is simplest — no options, just `NewCodec() encoding.Codec[[]byte, []byte]`.

**zstd, brotli** — same pattern as gzip, use their respective compression libraries.

- [ ] **Step 4: Update test files**

Tests change from decorator pattern to standalone + Pipe:

```go
// Old test pattern
c := gzip.NewCodec(json.NewCodec[Data]())
encoded, _ := c.Encode(Data{Name: "test"})

// New test pattern — test gzip standalone
c := gzip.NewCodec()
encoded, _ := c.Encode([]byte(`{"Name":"test"}`))
var decoded []byte
_ = c.Decode(encoded, &decoded)
// assert decoded == original bytes

// New test pattern — test with Pipe
combined := goencode.PipeCodec(json.NewCodec[Data](), gzip.NewCodec())
encoded, _ := combined.Encode(Data{Name: "test"})
var decoded Data
_ = combined.Decode(encoded, &decoded)
```

Update all `*_test.go` and `benchmark_test.go` files accordingly.

- [ ] **Step 5: Run tests**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./gzip/... ./flate/... ./snappy/... ./zstd/... ./brotli/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add gzip/ flate/ snappy/ zstd/ brotli/ && git commit -m "refactor: make compression codecs standalone Codec[[]byte,[]byte]

Remove decorator pattern. Each compression codec now operates on raw
bytes. Use goencode.PipeCodec() to compose with serialization codecs."
```

---

### Task 7: Migrate File Codec

File codec stays as decorator — wraps `Codec[T, []byte]` with own signature.

**Files:**
- Modify: `file/codec.go`
- Modify: `file/streamcodec.go`
- Modify: `file/option.go` (likely unchanged)
- Modify: `file/codec_test.go`, `file/streamcodec_test.go`

- [ ] **Step 1: Update `file/codec.go`**

Only change: the inner codec type from `encoding.Codec[T]` to `encoding.Codec[T, []byte]`.

```go
package file

import (
	"fmt"
	"os"
	"path/filepath"

	encoding "github.com/foomo/goencode"
)

// Codec encodes T to a file and decodes T from a file using an underlying Codec[T, []byte].
// Writes are atomic: data is written to a temporary file and renamed into place.
// It is safe for concurrent use.
type Codec[T any] struct {
	codec encoding.Codec[T, []byte]
	perm  os.FileMode
}

// NewCodec returns a file codec that delegates serialization to codec.
func NewCodec[T any](codec encoding.Codec[T, []byte], opts ...Option) *Codec[T] {
	o := options{
		perm: 0o644,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &Codec[T]{
		codec: codec,
		perm:  o.perm,
	}
}

// Encode serializes v and atomically writes the result to path.
func (c *Codec[T]) Encode(path string, v T) error {
	b, err := c.codec.Encode(v)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	f, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	tmp := f.Name()

	if _, err := f.Write(b); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Chmod(tmp, c.perm); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("setting file permissions: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// Decode reads the file at path and deserializes its contents into v.
func (c *Codec[T]) Decode(path string, v *T) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return c.codec.Decode(b, v)
}
```

- [ ] **Step 2: Update `file/streamcodec.go`**

Same change: inner codec type from `encoding.StreamCodec[T]` to `encoding.StreamCodec[T]` (StreamCodec signature is unchanged — still single param).

```go
package file

import (
	"fmt"
	"os"
	"path/filepath"

	encoding "github.com/foomo/goencode"
)

// StreamCodec encodes T to a file and decodes T from a file using an underlying StreamCodec[T].
// Writes are atomic: data is written to a temporary file and renamed into place.
// It is safe for concurrent use.
type StreamCodec[T any] struct {
	codec encoding.StreamCodec[T]
	perm  os.FileMode
}

// NewStreamCodec returns a file stream codec that delegates serialization to codec.
func NewStreamCodec[T any](codec encoding.StreamCodec[T], opts ...Option) *StreamCodec[T] {
	o := options{
		perm: 0o644,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &StreamCodec[T]{
		codec: codec,
		perm:  o.perm,
	}
}

// Encode serializes v and atomically writes the result to path.
func (c *StreamCodec[T]) Encode(path string, v T) error {
	dir := filepath.Dir(path)

	f, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	tmp := f.Name()

	if err := c.codec.Encode(f, v); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("encoding to temp file: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err := os.Chmod(tmp, c.perm); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("setting file permissions: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// Decode reads the file at path and deserializes its contents into v.
func (c *StreamCodec[T]) Decode(path string, v *T) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.codec.Decode(f, v)
}
```

- [ ] **Step 3: Update tests — constructor call stays same**

The file codec test should work with minimal changes since the API is the same. The only difference is `json.NewCodec[Data]()` now returns a struct instead of interface — but Go handles this transparently.

- [ ] **Step 4: Run tests**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./file/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add file/ && git commit -m "refactor: update file codec to accept Codec[T, []byte]"
```

---

### Task 8: Migrate Submodule Codecs (json/v2, yaml/v2, yaml/v3, yaml/v4, msgpack/*)

These are separate go.mod modules with external dependencies.

**Files:**
- Modify: `json/v2/codec.go`
- Modify: `yaml/v2/codec.go`, `yaml/v3/codec.go`, `yaml/v4/codec.go`
- Modify: `msgpack/tinylib/codec.go`, `msgpack/tinylib/streamcodec.go`
- Modify: `msgpack/vmihailenco/codec.go`, `msgpack/vmihailenco/streamcodec.go`
- Modify: all corresponding test and benchmark files

- [ ] **Step 1: Migrate `json/v2/codec.go`**

Same pattern as json/v1 but uses `github.com/go-json-experiment/json`. Return `encoding.Codec[T, []byte]`. If it has additional methods (`EncodeTo`/`DecodeFrom`), those can be dropped since StreamCodec covers streaming.

- [ ] **Step 2: Migrate yaml codecs**

All three yaml versions follow identical pattern — return `encoding.Codec[T, []byte]`.

- [ ] **Step 3: Migrate msgpack codecs**

**msgpack/tinylib** — has type constraints (`msgp.Marshaler`/`msgp.Unmarshaler`). The constraint stays but return type changes to `encoding.Codec[T, []byte]`. Keep the type assertion checks in the constructor.

**msgpack/vmihailenco** — standard pattern, return `encoding.Codec[T, []byte]`.

- [ ] **Step 4: Update go.mod files**

Each submodule's `go.mod` references the root module. Run `go mod tidy` in each:

```bash
for dir in json/v2 yaml/v2 yaml/v3 yaml/v4 msgpack/tinylib msgpack/vmihailenco; do
  (cd "$dir" && go mod tidy)
done
```

- [ ] **Step 5: Update tests, run all**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./json/v2/... ./yaml/v2/... ./yaml/v3/... ./yaml/v4/... ./msgpack/tinylib/... ./msgpack/vmihailenco/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add json/v2/ yaml/ msgpack/ && git commit -m "refactor: migrate submodule codecs (json/v2, yaml, msgpack) to Codec[S,T]"
```

---

### Task 9: Full Test Suite & Lint

**Files:** None new — validation only.

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make test`
Expected: PASS

- [ ] **Step 2: Run linter**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make lint`
Expected: PASS (or only pre-existing warnings)

- [ ] **Step 3: Fix any lint issues**

Address any new lint warnings introduced by the migration.

- [ ] **Step 4: Run race detector**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make test.race`
Expected: PASS

- [ ] **Step 5: Run full check**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make check`
Expected: PASS

- [ ] **Step 6: Commit any fixes**

```bash
git add -A && git commit -m "fix: address lint issues from codec interface migration"
```

(Skip if no fixes needed.)
