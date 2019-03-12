package binpack

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
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
		{nil, "0f"},
		{true, "04"},
		{false, "05"},
		{"", "20"},
		{"a", "2161"},
		{"hello", "2568656c6c6f"},
		{[]byte{}, "10"},
		{[]byte(nil), "10"},
		{[]byte("abc¢"), "15616263c2a2"},
		{[3]byte{1, 2, 3}, "13010203"},
		{[1]byte{}, "1100"},
		{[2]byte{1}, "120100"},
		{float32(3.14), "07c3f54840"},
		{float32(0), "0700000000"},
		{float32(-3.14), "07c3f548c0"},
		{float64(3.14), "061f85eb51b81e0940"},
		{float64(0), "060000000000000000"},
		{float64(-3.14), "061f85eb51b81e09c0"},
		{[]string{"a", "b", "c"}, "21612162216301"},
		{[3][2]int{}, "40400140400140400101"},
		{[2][3]string{}, "202020012020200101"},
		{[3]string{"a", "b", "c"}, "21612162216301"},
		{map[string]string(nil), "0301"},
		{map[int]string{1: "string"}, "4126737472696e6701"},
		{
			map[string]string{"a": "", "b": "", "c": "", "d": "", "e": ""},
			"21612021622021632021642021652001",
		},
		{int8(-1), "69"},
		{int32(1), "59"},
		{int64(math.MaxInt64), "ffffffffffffffffff40"},
		{uint8(8), "8848"},
		{uint64(math.MaxUint64), "ffffffffffffffffff41"},
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
