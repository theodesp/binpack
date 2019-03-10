package binpack

import (
	"io"
	"testing"
)

func TestDecBuffer_Drop(t *testing.T) {
	db := decBuffer{}
	db.Size(60)

	if db.Len() != 60 {
		t.Fatalf("decBuffer:Len expected length %v: got %v", 60, db.Len())
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("decBuffer:Drop should have panicked")
		}
	}()
	db.Drop(64)
	db.Drop(56)
	if db.Len() != 4 {
		t.Fatalf("decBuffer:Len expected length after drop %v: got %v", 4, db.Len())
	}
}

func TestDecBuffer_Reset(t *testing.T) {
	db := decBuffer{}
	db.Size(60)
	db.Reset()
	if db.Len() != 0 {
		t.Fatalf("decBuffer:Len expected length after Reset %v: got %v", 0, db.Len())
	}
}

func TestDecBuffer_Bytes(t *testing.T) {
	db := decBuffer{}
	db.Size(60)
	b := db.Bytes()
	if len(b) != 60 {
		t.Fatalf("decBuffer:Bytes expected return %v: bytes got %v", 60, db.Len())
	}
}

func TestDecBuffer_ReadByte(t *testing.T) {
	db := decBuffer{}
	_, err := db.ReadByte()
	if err != io.EOF {
		t.Fatalf("decBuffer:ReadByte expected EOF: got %v", err)
	}
	db.Size(1)
	b := db.Bytes()
	b[0] = byte('a')
	out, err := db.ReadByte()
	if err != nil {
		t.Fatalf("decBuffer:ReadByte expected no error: got %v", err)
	}

	if out != byte('a') {
		t.Fatalf("decBuffer:ReadByte expected %v: got %v", byte('a'), out)
	}
}

func TestDecBuffer_Read(t *testing.T) {
	db := decBuffer{}
	_, err := db.Read([]byte{})
	if err == io.EOF {
		t.Fatalf("decBuffer:Read expected nil: got %v", err)
	}
	_, err = db.Read([]byte{1})
	if err != io.EOF {
		t.Fatalf("decBuffer:Read expected EOF: got %v", err)
	}
	db.Size(4)
	b := db.Bytes()
	b[0] = byte('a')
	b[1] = byte('b')
	b[2] = byte('c')
	buf := make([]byte, 3)
	n, err := db.Read(buf)
	if err != nil {
		t.Fatalf("decBuffer:Read expected no error: got %v", err)
	}
	if n != 3 {
		t.Fatalf("decBuffer:Read expected to read %v bytes: got %v", 3, n)
	}
}
