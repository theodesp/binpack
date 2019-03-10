package binpack

import (
	"bytes"
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
		want []byte
	}{
		{
			nil,
			[]byte{byte(Nil)},
		},
		{
			true,
			[]byte{byte(True)},
		},
		{
			false,
			[]byte{byte(False)},
		},
	}
	var w bytes.Buffer
	enc := NewEncoder(&w)

	for _, test := range testCases {
		w.Reset()
		err := enc.Encode(test.in)
		if err != nil {
			t.Fatalf("binpack:Encode error %v", err)
		}
		got := w.Bytes()
		if !bytes.Equal(got, test.want) {
			t.Fatalf("%s != %s (in=%#v)", got, test.want, test.in)
		}
	}
}
