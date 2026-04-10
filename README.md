[![Build Status](https://github.com/foomo/goencode/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/goencode/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/goencode)](https://goreportcard.com/report/github.com/foomo/goencode)
[![GoDoc](https://godoc.org/github.com/foomo/goencode?status.svg)](https://godoc.org/github.com/foomo/goencode)

<p align="center">
  <img alt="goencode" src="docs/public/logo.png" width="400" height="400"/>
</p>

# goencode

> Generic encoding interfaces for Go with composable codecs.

## Features

- Generic `Codec[T]` and `StreamCodec[T]` interfaces with compile-time type safety
- Composable compression wrappers (gzip, flate, snappy, zstd, brotli) using the decorator pattern
- Stream support via `io.Reader`/`io.Writer` for memory-efficient pipelines
- Atomic file I/O with temp file + rename
- Zero dependencies in the core module

## Installation

```bash
go get github.com/foomo/goencode
```

## Core Interfaces

```go
// Byte-oriented
type Codec[T any] interface {
Encode(v T) ([]byte, error)
Decode(b []byte, v *T) error
}

// Stream-oriented
type StreamCodec[T any] interface {
Encode(w io.Writer, v T) error
Decode(r io.Reader, v *T) error
}
```

## Quick Start

```go
// Basic JSON encode/decode
c := json.NewCodec[User]()
b, err := c.Encode(User{Name: "Alice", Age: 30})
var u User
err = c.Decode(b, &u)

// Compose with compression
c := gzip.NewCodec[User](json.NewCodec[User]())

// Add atomic file persistence
fc := file.NewCodec[User](gzip.NewCodec[User](json.NewCodec[User]()))
err := fc.Encode("/tmp/user.json.gz", user)
```

## Available Codecs

| Category             | Packages                                                                                                                        |
|----------------------|---------------------------------------------------------------------------------------------------------------------------------|
| Serialization        | `json/v1`, `xml`, `gob`, `asn1`, `csv`, `pem`                                                                                   |
| Binary encoding      | `base64`, `base32`, `hex`, `ascii85`                                                                                            |
| Compression wrappers | `gzip`, `flate`, `snappy`\*, `zstd`\*, `brotli`\*                                                                               |
| Utility              | `file` (atomic read/write)                                                                                                      |
| Alternatives         | `json/v2`\* (go-json-experiment), `yaml/v2`\*, `yaml/v3`\*, `yaml/v4`\*, `toml`\*, `msgpack/tinylib`\*, `msgpack/vmihailenco`\* |

\* Submodule with separate `go.mod` — requires its own `go get`.

All codecs are safe for concurrent use.

<!-- BEGIN BENCHMARKS -->

## Benchmarks

> Measured with `go test -bench=. -benchmem` on arm64 (darwin).
> Results vary by hardware — use these as relative comparisons between codecs.

| Codec                 | Encode (ns/op) | Encode (B/op) | Encode (allocs/op) | Decode (ns/op) | Decode (B/op) | Decode (allocs/op) |
|-----------------------|---------------:|--------------:|-------------------:|---------------:|--------------:|-------------------:|
| `ascii85`             |          428.9 |           640 |                  1 |          870.7 |          4272 |                  4 |
| `asn1`                |           1192 |          1080 |                 18 |          563.4 |           552 |                  4 |
| `base32`              |          211.4 |           768 |                  1 |           1173 |          1248 |                  2 |
| `base64`              |          201.8 |           640 |                  1 |          172.0 |           480 |                  1 |
| `brotli`              |         121976 |       2165234 |                 27 |          10356 |         65288 |                 26 |
| `csv`                 |          420.2 |          4162 |                  2 |          715.6 |          4840 |                 22 |
| `file`                |         174240 |          1578 |                 12 |          49555 |          1760 |                 14 |
| `flate`               |          52005 |        816445 |                 20 |           8301 |         42648 |                 17 |
| `gob`                 |          965.1 |          2305 |                 18 |           5836 |          8320 |                167 |
| `gzip`                |          43324 |        816214 |                 20 |           8392 |         43352 |                 18 |
| `hex`                 |          308.3 |          1024 |                  1 |          309.1 |           480 |                  1 |
| `json/v1`             |          359.3 |           576 |                  1 |          467.7 |           552 |                  3 |
| `json/v2`             |          428.7 |           576 |                  1 |          456.9 |           552 |                  3 |
| `msgpack/tinylib`     |          91.76 |           704 |                  3 |          88.63 |           552 |                  4 |
| `msgpack/vmihailenco` |          187.2 |           688 |                  3 |          279.0 |           608 |                  6 |
| `pem`                 |          634.6 |          3312 |                 10 |           1448 |           580 |                  4 |
| `snappy`              |          734.2 |          1281 |                  2 |           2196 |          1352 |                 10 |
| `toml`                |           2276 |          6012 |                 47 |           5761 |          5400 |                 56 |
| `xml`                 |           1527 |          5131 |                  9 |           4868 |          4032 |                 57 |
| `yaml/v2`             |           7265 |          7200 |                 42 |           8109 |          9944 |                137 |
| `yaml/v3`             |           7507 |          9312 |                 42 |           9013 |         12256 |                151 |
| `yaml/v4`             |           7299 |          8640 |                 41 |           9155 |         12528 |                151 |
| `zstd`                |         350425 |      21432510 |                 57 |           8075 |         25092 |                 37 |

<!-- END BENCHMARKS -->

## How to Contribute

Contributions are welcome! Please read the [contributing guide](docs/CONTRIBUTING.md).

![Contributors](https://contributors-table.vercel.app/image?repo=foomo/goencode&width=50&columns=15)

## License

Distributed under MIT License, please see the [license](LICENSE) file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
