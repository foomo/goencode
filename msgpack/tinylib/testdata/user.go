package testdata

//go:generate msgp -o user_tinylib_gen.go -tests=false

import (
	"time"

	"github.com/foomo/goencode/internal/testdata"
)

// User mirrors User but carries tinylib/msgp codegen.
// Used to benchmark codegen-backed serialization against reflection-based User.
type User struct {
	Handle      string `msg:"handle"`
	Country     string `msg:"country"`
	Timestamp   int64  `msg:"timestamp"`
	Description string `msg:"description"`
}

func NewUserTinyLib() *User {
	return &User{
		Handle:      "@bench",
		Country:     "US",
		Timestamp:   time.Now().UnixNano(),
		Description: testdata.Text,
	}
}
