module github.com/foomo/goencode/msgpack/vmihailenco

go 1.26

replace github.com/foomo/goencode => ../../

require (
	github.com/foomo/goencode v0.0.0-00010101000000-000000000000
	github.com/vmihailenco/msgpack/v5 v5.4.1
)

require github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
