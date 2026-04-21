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
