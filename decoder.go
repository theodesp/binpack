package binpack

import (
	"bufio"
	"errors"
	"io"
	"reflect"
)

// A Decoder parses a decoded message and unpacks its values into the assigned variables.
// It is NOT safe for concurrent use by multiple
// goroutines.
type Decoder struct {
	r   io.Reader // source of the data
	buf decBuffer // buffer for more efficient i/o from r
	err error     // handle reader errors
}

// NewDecoder returns a new decoder that reads from the io.Reader.
// If r does not also implement io.ByteReader, it will be wrapped in a
// bufio.Reader.
func NewDecoder(r io.Reader) *Decoder {
	dec := new(Decoder)
	// We use the ability to read bytes as a plausible surrogate for buffering.
	if _, ok := r.(io.ByteReader); !ok {
		r = bufio.NewReader(r)
	}
	dec.r = r
	return dec
}

// Decode reads the next value from the input stream and stores
// it in the data represented by the empty interface value.
// If e is nil, the value will be discarded. Otherwise,
// the value underlying e must be a pointer to the
// correct type for the next data item received.
// If the input is at EOF, Decode returns io.EOF and
// does not modify e.
func (dec *Decoder) Decode(e interface{}) error {
	if e == nil {
		return dec.DecodeValue(reflect.Value{})
	}
	value := reflect.ValueOf(e)
	// If e represents a value as opposed to a pointer, the answer won't
	// get back to the caller. Make sure it's a pointer.
	if value.Type().Kind() != reflect.Ptr {
		dec.err = errors.New("binpack: attempt to decode into a non-pointer")
		return dec.err
	}
	return dec.DecodeValue(value)
}

// DecodeValue reads the next value from the input stream.
// If v is the zero reflect.Value (v.Kind() == Invalid), DecodeValue discards the value.
// Otherwise, it stores the value into v. In that case, v must represent
// a non-nil pointer to data or be an assignable reflect.Value (v.CanSet())
// If the input is at EOF, DecodeValue returns io.EOF and
// does not modify v.
func (dec *Decoder) DecodeValue(v reflect.Value) error {
	if v.IsValid() {
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			// That's okay, we'll store through the pointer.
		} else if !v.CanSet() {
			return errors.New("binpack: DecodeValue of unassignable value")
		}
	}

	dec.buf.Reset() // In case data lingers from previous invocation.
	dec.err = nil
	id := dec.decodeType(false)
	if dec.err == nil {
		dec.decode(id, v)
	}
	return dec.err
}

// decodeType parses
// and returns the type id of the next value. It returns 0 at
// EOF.  Upon return, the remainder of dec.buf is the value to be
// decoded. If this is an interface value, it can be ignored by
// resetting that buffer.
func (dec *Decoder) decodeType(isInterface bool) Code {
	return 0
}

// decode decodes the data stream representing a value and stores it in value.
func (dec *Decoder) decode(_ Code, _ reflect.Value) {
	defer catchError(&dec.err)
}
