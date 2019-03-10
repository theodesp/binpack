package binpack

import (
	"bytes"
	"testing"
)

func TestEncBuffer_WriteByte(t *testing.T) {
	eb := encBuffer{}
	eb.WriteByte(byte('a'))
	out := eb.Bytes()
	in := []byte{'a'}
	if !bytes.Equal(out, in) {
		t.Fatalf("encBuffer:WriteByte got %v; wanted %v", out, in)
	}
}

func TestEncBuffer_WriteString(t *testing.T) {
	eb := encBuffer{}
	eb.WriteString("abc")
	out := eb.Bytes()
	in := []byte{'a', 'b', 'c'}
	if !bytes.Equal(out, in) {
		t.Fatalf("encBuffer:WriteString got %v; wanted %v", out, in)
	}
}

func TestEncBuffer_Write(t *testing.T) {
	eb := encBuffer{}
	n, _ := eb.Write([]byte{'a', 'b', 'c'})
	if n != 3 {
		t.Fatalf("encBuffer:Write wrote %v bytes; wanted %v bytes", n, 3)
	}
	out := eb.Bytes()
	in := []byte{'a', 'b', 'c'}
	if !bytes.Equal(out, in) {
		t.Fatalf("encBuffer:Write got %v; wanted %v", out, in)
	}
}

func TestEncBuffer_Len(t *testing.T) {
	eb := encBuffer{}
	_, err := eb.Write([]byte("abc日"))
	if err != nil {
		t.Fatal("encBuffer:Write error", err)
	}
	if eb.Len() != 6 {
		t.Fatalf("encBuffer:Write wrote %v bytes; wanted %v bytes", eb.Len(), 6)
	}
}

func TestEncBuffer_Reset(t *testing.T) {
	eb := encBuffer{}
	_, err := eb.Write([]byte("abc日"))
	if err != nil {
		t.Fatal("encBuffer:Write error", err)
	}
	eb.Reset()
	out := eb.Bytes()
	if !bytes.Equal(out, []byte{}) {
		t.Fatalf("encBuffer:Reset got %v; wanted %v", out, []byte{})
	}
}
