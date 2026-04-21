# Standalone Encoder/Decoder Exports Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Export standalone `Encoder` and `Decoder` functions from all subpackages so consumer APIs can depend on a single direction without requiring a full `Codec`.

**Architecture:** Each subpackage gets bare `Encoder`/`Decoder` funcs (simple codecs) or `NewEncoder`/`NewDecoder` constructors (configurable codecs). Existing `NewCodec` constructors are refactored to delegate to these new exports. The reference implementation is `json/v1/codec.go` which already has this pattern.

**Tech Stack:** Go generics, `go.work` multi-module workspace

---

### Task 1: Simple serialization codecs — stdlib (xml, gob, asn1, csv)

These packages use stdlib only, live in the root module, and follow the same pattern: extract inline encode/decode closures into named generic funcs.

**Files:**
- Modify: `xml/codec.go`
- Modify: `gob/codec.go`
- Modify: `asn1/codec.go`
- Modify: `csv/codec.go`

- [ ] **Step 1: Refactor `xml/codec.go`**

```go
package xml

import (
	stdxml "encoding/xml"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes T to XML bytes.
func Encoder[T any](v T) ([]byte, error) {
	return stdxml.Marshal(v)
}

// Decoder decodes XML bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return stdxml.Unmarshal(b, v)
}

// NewCodec returns an XML codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 2: Refactor `gob/codec.go`**

```go
package gob

import (
	"bytes"
	stdgob "encoding/gob"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// Encoder encodes T to gob bytes.
func Encoder[T any](v T) ([]byte, error) {
	buf := sync.Get()
	defer sync.Put(buf)

	if err := stdgob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

// Decoder decodes gob bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return stdgob.NewDecoder(bytes.NewReader(b)).Decode(v)
}

// NewCodec returns a gob codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 3: Refactor `asn1/codec.go`**

```go
package asn1

import (
	stdasn1 "encoding/asn1"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes T to ASN.1 bytes.
func Encoder[T any](v T) ([]byte, error) {
	return stdasn1.Marshal(v)
}

// Decoder decodes ASN.1 bytes into T.
func Decoder[T any](b []byte, v *T) error {
	_, err := stdasn1.Unmarshal(b, v)
	return err
}

// NewCodec returns an ASN1 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 4: Refactor `csv/codec.go`**

```go
package csv

import (
	"bytes"
	stdcsv "encoding/csv"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// Encoder encodes [][]string to CSV bytes.
func Encoder(v [][]string) ([]byte, error) {
	buf := sync.Get()
	defer sync.Put(buf)

	cw := stdcsv.NewWriter(buf)
	if err := cw.WriteAll(v); err != nil {
		return nil, err
	}

	cw.Flush()

	if err := cw.Error(); err != nil {
		return nil, err
	}

	return append([]byte(nil), buf.Bytes()...), nil
}

// Decoder decodes CSV bytes into [][]string.
func Decoder(b []byte, v *[][]string) error {
	records, err := stdcsv.NewReader(bytes.NewReader(b)).ReadAll()
	if err != nil {
		return err
	}

	*v = records

	return nil
}

// NewCodec returns a CSV codec for [][]string.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[][]string, []byte] {
	return encoding.Codec[[][]string, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 5: Run tests for stdlib codecs**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./xml/... ./gob/... ./asn1/... ./csv/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add xml/codec.go gob/codec.go asn1/codec.go csv/codec.go
git commit -m "refactor: export standalone Encoder/Decoder in xml, gob, asn1, csv"
```

---

### Task 2: Simple encoding codecs (base64, base32, hex, ascii85, pem)

These are `[]byte ↔ []byte` (or `*pem.Block ↔ []byte`) — no type parameters.

**Files:**
- Modify: `base64/codec.go`
- Modify: `base32/codec.go`
- Modify: `hex/codec.go`
- Modify: `ascii85/codec.go`
- Modify: `pem/codec.go`

- [ ] **Step 1: Refactor `base64/codec.go`**

```go
package base64

import (
	stdbase64 "encoding/base64"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to Base64.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdbase64.StdEncoding.EncodedLen(len(v)))
	stdbase64.StdEncoding.Encode(dst, v)

	return dst, nil
}

// Decoder decodes Base64 bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdbase64.StdEncoding.DecodedLen(len(b)))

	n, err := stdbase64.StdEncoding.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Base64 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 2: Refactor `base32/codec.go`**

```go
package base32

import (
	stdbase32 "encoding/base32"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to Base32.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdbase32.StdEncoding.EncodedLen(len(v)))
	stdbase32.StdEncoding.Encode(dst, v)

	return dst, nil
}

// Decoder decodes Base32 bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdbase32.StdEncoding.DecodedLen(len(b)))

	n, err := stdbase32.StdEncoding.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Base32 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 3: Refactor `hex/codec.go`**

```go
package hex

import (
	stdhex "encoding/hex"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to hexadecimal.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdhex.EncodedLen(len(v)))
	stdhex.Encode(dst, v)

	return dst, nil
}

// Decoder decodes hexadecimal bytes.
func Decoder(b []byte, v *[]byte) error {
	dst := make([]byte, stdhex.DecodedLen(len(b)))

	n, err := stdhex.Decode(dst, b)
	if err != nil {
		return err
	}

	*v = dst[:n]

	return nil
}

// NewCodec returns a Hex codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 4: Refactor `ascii85/codec.go`**

```go
package ascii85

import (
	"bytes"
	stdascii85 "encoding/ascii85"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes bytes to ASCII85.
func Encoder(v []byte) ([]byte, error) {
	dst := make([]byte, stdascii85.MaxEncodedLen(len(v)))
	n := stdascii85.Encode(dst, v)

	return dst[:n], nil
}

// Decoder decodes ASCII85 bytes.
func Decoder(b []byte, v *[]byte) error {
	buf := bytes.NewBuffer(make([]byte, 0, len(b)))

	r := stdascii85.NewDecoder(bytes.NewReader(b))
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	*v = buf.Bytes()

	return nil
}

// NewCodec returns an ASCII85 codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 5: Refactor `pem/codec.go`**

```go
package pem

import (
	stdpem "encoding/pem"
	"errors"

	encoding "github.com/foomo/goencode"
)

// Encoder encodes a PEM block to bytes.
func Encoder(v *stdpem.Block) ([]byte, error) {
	return stdpem.EncodeToMemory(v), nil
}

// Decoder decodes bytes into a PEM block.
func Decoder(b []byte, v **stdpem.Block) error {
	block, _ := stdpem.Decode(b)
	if block == nil {
		return errors.New("pem: no PEM block found")
	}

	*v = block

	return nil
}

// NewCodec returns a PEM codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[*stdpem.Block, []byte] {
	return encoding.Codec[*stdpem.Block, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 6: Run tests for encoding codecs**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./base64/... ./base32/... ./hex/... ./ascii85/... ./pem/...`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add base64/codec.go base32/codec.go hex/codec.go ascii85/codec.go pem/codec.go
git commit -m "refactor: export standalone Encoder/Decoder in base64, base32, hex, ascii85, pem"
```

---

### Task 3: Simple compression codec (snappy)

Snappy takes no options, so it gets bare funcs.

**Files:**
- Modify: `snappy/codec.go`

- [ ] **Step 1: Refactor `snappy/codec.go`**

```go
package snappy

import (
	encoding "github.com/foomo/goencode"
	"github.com/golang/snappy"
)

// Encoder compresses bytes using Snappy.
func Encoder(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

// Decoder decompresses Snappy bytes.
func Decoder(data []byte, v *[]byte) error {
	decoded, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	*v = decoded

	return nil
}

// NewCodec returns a Snappy compression codec.
// It is safe for concurrent use.
func NewCodec() encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: Encoder,
		Decode: Decoder,
	}
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode/snappy && go test -tags=safe ./...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add snappy/codec.go
git commit -m "refactor: export standalone Encoder/Decoder in snappy"
```

---

### Task 4: Configurable compression codecs (gzip, flate, zstd, brotli)

These accept `Option` variadic args. They get `NewEncoder`/`NewDecoder` constructors. The closures inside `NewCodec` are extracted into these constructors, and `NewCodec` delegates to them.

**Files:**
- Modify: `gzip/codec.go`
- Modify: `flate/codec.go`
- Modify: `zstd/codec.go`
- Modify: `brotli/codec.go`

- [ ] **Step 1: Refactor `gzip/codec.go`**

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

// NewEncoder returns a gzip compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: gzip.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
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
	}
}

// NewDecoder returns a gzip decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
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
	}
}

// NewCodec returns a gzip compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
```

- [ ] **Step 2: Refactor `flate/codec.go`**

```go
package flate

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewEncoder returns a DEFLATE compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: flate.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
		buf := sync.Get()
		defer sync.Put(buf)

		w, err := flate.NewWriter(buf, o.level)
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
	}
}

// NewDecoder returns a DEFLATE decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
		r := flate.NewReader(bytes.NewReader(data))
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
			return fmt.Errorf("flate: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
		}

		*v = decoded

		return nil
	}
}

// NewCodec returns a DEFLATE compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
```

- [ ] **Step 3: Refactor `zstd/codec.go`**

```go
package zstd

import (
	encoding "github.com/foomo/goencode"
	"github.com/klauspost/compress/zstd"
)

// NewEncoder returns a Zstandard compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: zstd.SpeedDefault,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
		enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(o.level))
		if err != nil {
			return nil, err
		}
		defer enc.Close()

		return enc.EncodeAll(data, nil), nil
	}
}

// NewDecoder returns a Zstandard decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
		dopts := []zstd.DOption{}
		if o.maxDecodedSize > 0 {
			dopts = append(dopts, zstd.WithDecoderMaxMemory(uint64(o.maxDecodedSize)))
		}

		dec, err := zstd.NewReader(nil, dopts...)
		if err != nil {
			return err
		}
		defer dec.Close()

		decoded, err := dec.DecodeAll(data, nil)
		if err != nil {
			return err
		}

		*v = decoded

		return nil
	}
}

// NewCodec returns a Zstandard compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
```

- [ ] **Step 4: Refactor `brotli/codec.go`**

```go
package brotli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	encoding "github.com/foomo/goencode"
	"github.com/foomo/goencode/internal/sync"
)

// NewEncoder returns a Brotli compression encoder.
func NewEncoder(opts ...Option) encoding.Encoder[[]byte, []byte] {
	o := options{
		level: brotli.DefaultCompression,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte) ([]byte, error) {
		buf := sync.Get()
		defer sync.Put(buf)

		w := brotli.NewWriterLevel(buf, o.level)

		if _, err := w.Write(data); err != nil {
			return nil, err
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		return append([]byte(nil), buf.Bytes()...), nil
	}
}

// NewDecoder returns a Brotli decompression decoder.
func NewDecoder(opts ...Option) encoding.Decoder[[]byte, []byte] {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(data []byte, v *[]byte) error {
		r := brotli.NewReader(bytes.NewReader(data))

		var src io.Reader = r
		if o.maxDecodedSize > 0 {
			src = io.LimitReader(r, o.maxDecodedSize+1)
		}

		decoded, err := io.ReadAll(src)
		if err != nil {
			return err
		}

		if o.maxDecodedSize > 0 && int64(len(decoded)) > o.maxDecodedSize {
			return fmt.Errorf("brotli: decompressed size exceeds limit of %d bytes", o.maxDecodedSize)
		}

		*v = decoded

		return nil
	}
}

// NewCodec returns a Brotli compression codec.
// It is safe for concurrent use.
func NewCodec(opts ...Option) encoding.Codec[[]byte, []byte] {
	return encoding.Codec[[]byte, []byte]{
		Encode: NewEncoder(opts...),
		Decode: NewDecoder(opts...),
	}
}
```

- [ ] **Step 5: Run tests for compression codecs**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && go test -tags=safe ./gzip/... ./flate/...`
Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode/zstd && go test -tags=safe ./...`
Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode/brotli && go test -tags=safe ./...`
Expected: PASS for all

- [ ] **Step 6: Commit**

```bash
git add gzip/codec.go flate/codec.go zstd/codec.go brotli/codec.go
git commit -m "refactor: export NewEncoder/NewDecoder in gzip, flate, zstd, brotli"
```

---

### Task 5: Submodule serialization codecs (json/v2, yaml/v2, yaml/v3, yaml/v4, msgpack/tinylib, msgpack/vmihailenco, toml)

These live in separate go.mod submodules. Same bare-func pattern as Task 1.

**Files:**
- Modify: `json/v2/codec.go`
- Modify: `yaml/v2/codec.go`
- Modify: `yaml/v3/codec.go`
- Modify: `yaml/v4/codec.go`
- Modify: `msgpack/tinylib/codec.go`
- Modify: `msgpack/vmihailenco/codec.go`
- Modify: `toml/codec.go`

- [ ] **Step 1: Refactor `toml/codec.go`**

```go
package toml

import (
	encoding "github.com/foomo/goencode"

	"github.com/BurntSushi/toml"
)

// Encoder encodes T to TOML bytes.
func Encoder[T any](v T) ([]byte, error) {
	return toml.Marshal(v)
}

// Decoder decodes TOML bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return toml.Unmarshal(b, v)
}

// NewCodec returns a TOML codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 2: Refactor `json/v2/codec.go`**

Note: `json/v2` uses `github.com/go-json-experiment/json`. Read the current file to get the exact encode/decode logic before extracting. The current `NewCodec` uses `json.Marshal`/`json.Unmarshal` from that package.

```go
package json

import (
	encoding "github.com/foomo/goencode"
	"github.com/go-json-experiment/json"
)

// Encoder encodes T to JSON bytes (v2).
func Encoder[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

// Decoder decodes JSON bytes into T (v2).
func Decoder[T any](b []byte, v *T) error {
	return json.Unmarshal(b, v)
}

// NewCodec returns a JSON v2 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

Preserve any existing `NewStreamCodec` function unchanged at the bottom of the file.

- [ ] **Step 3: Refactor `yaml/v2/codec.go`**

```go
package yaml

import (
	encoding "github.com/foomo/goencode"
	"go.yaml.in/yaml/v2"
)

// Encoder encodes T to YAML v2 bytes.
func Encoder[T any](v T) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decoder decodes YAML v2 bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return yaml.Unmarshal(b, v)
}

// NewCodec returns a YAML v2 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 4: Refactor `yaml/v3/codec.go`**

```go
package yaml

import (
	encoding "github.com/foomo/goencode"
	"gopkg.in/yaml.v3"
)

// Encoder encodes T to YAML v3 bytes.
func Encoder[T any](v T) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decoder decodes YAML v3 bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return yaml.Unmarshal(b, v)
}

// NewCodec returns a YAML v3 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 5: Refactor `yaml/v4/codec.go`**

```go
package yaml

import (
	encoding "github.com/foomo/goencode"
	"github.com/goccy/go-yaml"
)

// Encoder encodes T to YAML bytes.
func Encoder[T any](v T) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decoder decodes YAML bytes into T.
func Decoder[T any](b []byte, v *T) error {
	return yaml.Unmarshal(b, v)
}

// NewCodec returns a YAML v4 codec for T.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 6: Refactor `msgpack/vmihailenco/codec.go`**

```go
package msgpack

import (
	encoding "github.com/foomo/goencode"
	"github.com/vmihailenco/msgpack/v5"
)

// Encoder encodes T to msgpack bytes (vmihailenco).
func Encoder[T any](v T) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decoder decodes msgpack bytes into T (vmihailenco).
func Decoder[T any](b []byte, v *T) error {
	return msgpack.Unmarshal(b, v)
}

// NewCodec returns a msgpack codec for T backed by vmihailenco/msgpack.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 7: Refactor `msgpack/tinylib/codec.go`**

This one is special — it checks for `msgp.Marshaler`/`msgp.Unmarshaler` interfaces. The `Encoder` and `Decoder` funcs must retain this runtime check.

```go
package msgpack

import (
	"fmt"

	encoding "github.com/foomo/goencode"
	"github.com/tinylib/msgp/msgp"
)

// Encoder encodes T to msgpack bytes (tinylib).
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Marshaler.
func Encoder[T any](v T) ([]byte, error) {
	if m, ok := any(v).(msgp.Marshaler); ok {
		return m.MarshalMsg(nil)
	}

	if m, ok := any(&v).(msgp.Marshaler); ok {
		return m.MarshalMsg(nil)
	}

	return nil, fmt.Errorf("msgpack: %T does not implement msgp.Marshaler", v)
}

// Decoder decodes msgpack bytes into T (tinylib).
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Unmarshaler.
func Decoder[T any](b []byte, v *T) error {
	if u, ok := any(v).(msgp.Unmarshaler); ok {
		_, err := u.UnmarshalMsg(b)
		return err
	}

	return fmt.Errorf("msgpack: %T does not implement msgp.Unmarshaler", v)
}

// NewCodec returns a msgpack codec for T backed by tinylib/msgp.
// T must have msgp code generation (go:generate msgp) so that
// *T implements msgp.Marshaler and msgp.Unmarshaler.
// It is safe for concurrent use.
func NewCodec[T any]() encoding.Codec[T, []byte] {
	return encoding.Codec[T, []byte]{
		Encode: Encoder[T],
		Decode: Decoder[T],
	}
}
```

- [ ] **Step 8: Run tests for all submodule codecs**

Run each in its own module directory:
```bash
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/toml && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/json/v2 && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/yaml/v2 && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/yaml/v3 && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/yaml/v4 && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/msgpack/vmihailenco && go test -tags=safe ./...
cd /Users/franklin/Workingcopies/github.com/foomo/goencode/msgpack/tinylib && go test -tags=safe ./...
```
Expected: PASS for all

- [ ] **Step 9: Commit**

```bash
git add toml/codec.go json/v2/codec.go yaml/v2/codec.go yaml/v3/codec.go yaml/v4/codec.go msgpack/tinylib/codec.go msgpack/vmihailenco/codec.go
git commit -m "refactor: export standalone Encoder/Decoder in submodule codecs"
```

---

### Task 6: Run full CI check

- [ ] **Step 1: Run lint**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make lint`
Expected: PASS (no new lint issues)

- [ ] **Step 2: Run full test suite**

Run: `cd /Users/franklin/Workingcopies/github.com/foomo/goencode && make test`
Expected: PASS

- [ ] **Step 3: Fix any issues found by lint or tests**

If lint reports issues (e.g. unused imports after refactor), fix them and re-run.

- [ ] **Step 4: Final commit if fixes were needed**

```bash
git add -u
git commit -m "fix: address lint issues from Encoder/Decoder export refactor"
```
