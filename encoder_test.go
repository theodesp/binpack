package binpack

import (
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
