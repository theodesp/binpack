package binpack

// tooBig provides a sanity check for sizes; used in several places. Upper limit
// of is 1GB on 32-bit systems, 8GB on 64-bit, allowing room to grow a little
// without overflow.
const tooBig = (1 << 30) << (^uint(0) >> 62)

// encBuffer is an extremely simple, fast implementation of a write-only byte buffer.
// It never returns a non-nil error, but Write returns an error value so it matches io.Writer.
type encBuffer struct {
	data []byte
	buf  [64]byte
}

func (e *encBuffer) WriteByte(c byte) {
	e.data = append(e.data, c)
}

func (e *encBuffer) WriteCode(c Code) {
	e.data = append(e.data, byte(c))
}

func (e *encBuffer) Write(p []byte) (int, error) {
	e.data = append(e.data, p...)
	return len(p), nil
}

func (e *encBuffer) WriteString(s string) {
	e.data = append(e.data, s...)
}

func (e *encBuffer) Len() int {
	return len(e.data)
}

func (e *encBuffer) Bytes() []byte {
	return e.data
}

func (e *encBuffer) Reset() {
	if len(e.data) >= tooBig {
		e.data = e.buf[0:0]
	} else {
		e.data = e.data[0:0]
	}
}
