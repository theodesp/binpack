package binpack

import (
	"io"
	"reflect"
)

// An Encoder manages the transmission of type and data information to the
// other side of a connection. It is NOT safe for concurrent use by multiple
// goroutines.
type Encoder struct {
	w   io.Writer // the writer to write to
	buf encBuffer // buffer to use when encoding data
	err error
}

// NewEncoder returns a new encoder that will transmit on the io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	enc := new(Encoder)
	enc.w = w
	return enc
}

// Encode transmits the data item represented by the empty interface value
func (enc *Encoder) Encode(e interface{}) error {
	return enc.EncodeValue(reflect.ValueOf(e))
}

// EncodeValue transmits the data item represented by the reflection value,
func (enc *Encoder) EncodeValue(value reflect.Value) error {
	enc.err = nil
	enc.buf.Reset()
	// Encode the object.
	enc.encode(value)
	if enc.err == nil {
		enc.writeTo(enc.w)
	}
	return enc.err
}

// encode nil into one byte to buffer.
//
// +-----------+
// | 0000 1111 |   0x0f
// +-----------+
func (enc *Encoder) encodeNil() {
	enc.buf.WriteCode(Nil)
}

// writeTo sends the data item to the writer
func (enc *Encoder) writeTo(w io.Writer) {
	// Write the data.
	_, err := w.Write(enc.buf.Bytes())
	// Drain the buffer and restore the space.
	enc.buf.Reset()
	if err != nil {
		enc.setError(err)
	}
}

func (enc *Encoder) setError(err error) {
	if enc.err == nil { // remember the first.
		enc.err = err
	}
}

func (enc *Encoder) encode(v reflect.Value) {
	defer catchError(&enc.err)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		panic("binpack: cannot encode nil pointer of type " + v.Type().String())
	}
	switch v.Kind() {
	case reflect.Invalid: // nil
		enc.encodeNil()
	default:
	}
}
