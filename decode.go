package binpack

import "io"

// decBuffer is an extremely simple, fast implementation of a read-only byte buffer.
// It is initialized by calling Size and then copying the data into the slice returned by Bytes().
type decBuffer struct {
	data   []byte
	offset int // Read offset.
}

func (d *decBuffer) Read(p []byte) (int, error) {
	n := copy(p, d.data[d.offset:])
	if n == 0 && len(p) != 0 {
		return 0, io.EOF
	}
	d.offset += n
	return n, nil
}

func (d *decBuffer) Drop(n int) {
	if n > d.Len() {
		panic("drop")
	}
	d.offset += n
}

// Size grows the buffer to exactly n bytes, so d.Bytes() will
// return a slice of length n. Existing data is first discarded.
func (d *decBuffer) Size(n int) {
	d.Reset()
	if cap(d.data) < n {
		d.data = make([]byte, n)
	} else {
		d.data = d.data[0:n]
	}
}

func (d *decBuffer) ReadByte() (byte, error) {
	if d.offset >= len(d.data) {
		return 0, io.EOF
	}
	c := d.data[d.offset]
	d.offset++
	return c, nil
}

func (d *decBuffer) Len() int {
	return len(d.data) - d.offset
}

func (d *decBuffer) Bytes() []byte {
	return d.data[d.offset:]
}

func (d *decBuffer) Reset() {
	d.data = d.data[0:0]
	d.offset = 0
}
