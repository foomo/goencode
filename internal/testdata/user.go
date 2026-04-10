package testdata

import "time"

type User struct {
	Handle      string `json:"handle"`
	Country     string `json:"country"`
	Timestamp   int64  `json:"timestamp"`
	Description string `json:"description"`
}

func NewUser() *User {
	return &User{
		Handle:      "@bench",
		Country:     "US",
		Timestamp:   time.Now().UnixNano(),
		Description: Text,
	}
}
