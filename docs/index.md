---
layout: home

hero:
  name: goencode
  text: Generic Encoding for Go
  tagline: Composable, type-safe codec interfaces. Serialize, compress, and persist data with a single API.
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
    details: Codec[T] and StreamCodec[T] use Go generics so encode/decode operations are statically typed at compile time.
  - title: Composable Wrappers
    details: Layer gzip, flate, snappy, or zstd compression on any codec with a single function call using the decorator pattern.
  - title: Streaming Support
    details: StreamCodec[T] reads and writes directly to io.Reader/io.Writer for memory-efficient pipelines and network I/O.
  - title: Atomic File I/O
    details: The file codec writes to a temp file and renames into place, preventing partial writes and data corruption.
  - title: Zero Dependencies
    details: The core module uses only the Go standard library. Submodules for yaml, json2, snappy, and zstd are separate imports.
  - title: Concurrent Safe
    details: All codecs are safe for concurrent use. Stateless serializers share nothing between goroutines.
---
