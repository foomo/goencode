# Examples

Complete, copy-pasteable examples. Each is a standalone `main` function.

## JSON Round-Trip

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/json/v1"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    c := json.NewCodec[User]()

    b, err := c.Encode(User{Name: "Alice", Age: 30}) // [!code highlight]
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(b)) // {"name":"Alice","age":30}

    var u User
    if err := c.Decode(b, &u); err != nil { // [!code highlight]
        log.Fatal(err)
    }
    fmt.Printf("Name: %s, Age: %d\n", u.Name, u.Age)
}
```

## XML with Struct Tags

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/xml"
)

type Data struct {
    XMLName struct{} `xml:"data"`
    Name    string   `xml:"name"`
}

func main() {
    c := xml.NewCodec[Data]()

    b, err := c.Encode(Data{Name: "example-123"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(b)) // <data><name>example-123</name></data>

    var d Data
    if err := c.Decode(b, &d); err != nil {
        log.Fatal(err)
    }
    fmt.Println(d.Name) // example-123
}
```

## CSV Records

```go
package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/foomo/goencode/csv"
)

func main() {
    c := csv.NewStreamCodec()

    records := [][]string{
        {"name", "age"},
        {"Alice", "30"},
    }

    var buf bytes.Buffer
    if err := c.Encode(&buf, records); err != nil { // [!code highlight]
        log.Fatal(err)
    }

    var decoded [][]string
    if err := c.Decode(&buf, &decoded); err != nil { // [!code highlight]
        log.Fatal(err)
    }
    fmt.Println(decoded) // [[name age] [Alice 30]]
}
```

## PEM Block Encoding

```go
package main

import (
    "bytes"
    stdpem "encoding/pem"
    "fmt"
    "log"

    "github.com/foomo/goencode/pem"
)

func main() {
    c := pem.NewStreamCodec()

    block := &stdpem.Block{
        Type:  "TEST",
        Bytes: []byte("hello"),
    }

    var buf bytes.Buffer
    if err := c.Encode(&buf, block); err != nil {
        log.Fatal(err)
    }

    var decoded *stdpem.Block
    if err := c.Decode(&buf, &decoded); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Type: %s, Data: %s\n", decoded.Type, string(decoded.Bytes))
    // Type: TEST, Data: hello
}
```

## Base64 Encoding Bytes

```go
package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/foomo/goencode/base64"
)

func main() {
    c := base64.NewStreamCodec()

    var buf bytes.Buffer
    if err := c.Encode(&buf, []byte("hello")); err != nil {
        log.Fatal(err)
    }
    fmt.Println(buf.String()) // aGVsbG8= // [!code highlight]

    var decoded []byte
    if err := c.Decode(&buf, &decoded); err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(decoded)) // hello
}
```

## Compressed JSON (gzip)

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type Event struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
}

func main() {
    c := gzip.NewCodec[Event](json.NewCodec[Event]()) // [!code highlight]

    b, err := c.Encode(Event{Type: "click", Payload: "button-1"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Compressed size: %d bytes\n", len(b))

    var e Event
    if err := c.Decode(b, &e); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Type: %s, Payload: %s\n", e.Type, e.Payload)
}
```

## Compressed JSON (brotli)

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/brotli"
    "github.com/foomo/goencode/json/v1"
)

type Event struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
}

func main() {
    c := brotli.NewCodec[Event](json.NewCodec[Event]()) // [!code highlight]

    b, err := c.Encode(Event{Type: "click", Payload: "button-1"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Compressed size: %d bytes\n", len(b))

    var e Event
    if err := c.Decode(b, &e); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Type: %s, Payload: %s\n", e.Type, e.Payload)
}
```

## TOML Round-Trip

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/toml"
)

type Config struct {
    Host string
    Port int
}

func main() {
    c := toml.NewCodec[Config]()

    b, err := c.Encode(Config{Host: "localhost", Port: 8080}) // [!code highlight]
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(b))

    var cfg Config
    if err := c.Decode(b, &cfg); err != nil { // [!code highlight]
        log.Fatal(err)
    }
    fmt.Printf("Host: %s, Port: %d\n", cfg.Host, cfg.Port)
}
```

## Atomic File Persistence

```go
package main

import (
    "fmt"
    "log"

    "github.com/foomo/goencode/file"
    "github.com/foomo/goencode/gzip"
    "github.com/foomo/goencode/json/v1"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

func main() {
    fc := file.NewCodec[Config]( // [!code highlight]
        gzip.NewCodec[Config](json.NewCodec[Config]()), // [!code highlight]
        file.WithPermissions(0o600), // [!code highlight]
    ) // [!code highlight]

    cfg := Config{Host: "localhost", Port: 8080}

    // Write atomically: JSON → gzip → temp file → rename
    if err := fc.Encode("/tmp/config.json.gz", cfg); err != nil {
        log.Fatal(err)
    }

    // Read back: file → gunzip → JSON
    var loaded Config
    if err := fc.Decode("/tmp/config.json.gz", &loaded); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Host: %s, Port: %d\n", loaded.Host, loaded.Port)
}
```
