package binpack

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"
)

type errorWriter struct{}

func (ew errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("forced error")
}

func TestWriter_ErrorOnWrite(t *testing.T) {
	w := errorWriter{}
	err := NewEncoder(w).Encode(1)
	if err == nil {
		t.Fatal("expected error from writer got none")
	}
}

func TestWriter_EncodeNilPointer(t *testing.T) {
	var w bytes.Buffer
	var p *interface{}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("encode should have panicked on nil interface pointer")
		}
	}()
	_ = NewEncoder(&w).Encode(p)
}

func TestEncoder(t *testing.T) {
	testCases := []struct {
		in   interface{}
		want string
	}{
		{
			nil,
			"0f",
		},
		{
			true,
			"04",
		},
		{
			false,
			"05",
		},
		{"", "20"},
		{"a", "2161"},
		{"hello", "2568656c6c6f"},
	}
	var w bytes.Buffer
	enc := NewEncoder(&w)

	for _, test := range testCases {
		w.Reset()
		err := enc.Encode(test.in)
		if err != nil {
			t.Fatalf("binpack:Encode error %v", err)
		}
		got := hex.EncodeToString(w.Bytes())
		if got != test.want {
			t.Fatalf("%s != %s (in=%#v)", got, test.want, test.in)
		}
	}
}
