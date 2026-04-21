---
layout: home

hero:
  name: goencode
  text: Generic Encoding for Go
  tagline: Composable, type-safe codec function types. Serialize, compress, and persist data with a single API.
  image:
    src: /logo.png
    alt: goencode
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View Codecs
      link: /guide/codecs

features:
  - title: Type-Safe Generics
    details: Codec[S, T] and StreamCodec[S] use Go generics so encode/decode operations are statically typed at compile time.
  - title: Composable Pipelines
    details: Chain any two codecs with PipeCodec — e.g., JSON → gzip, JSON → base64 — with full type safety at compile time.
  - title: Streaming Support
    details: StreamCodec[S] reads and writes directly to io.Reader/io.Writer for memory-efficient pipelines and network I/O.
  - title: Atomic File I/O
    details: The file codec writes to a temp file and renames into place, preventing partial writes and data corruption.
  - title: Zero Dependencies
    details: The core module uses only the Go standard library. Submodules for yaml, json2, snappy, and zstd are separate imports.
  - title: Concurrent Safe
    details: All codecs are safe for concurrent use. Stateless serializers share nothing between goroutines.
---
