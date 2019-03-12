package binpack

import (
	"io"
	"math"
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
	case reflect.Bool:
		enc.encodeBool(v.Bool())
	case reflect.Float32:
		enc.encodeFloat32(float32(v.Float()))
	case reflect.Float64:
		enc.encodeFloat64(v.Float())
	case reflect.String:
		enc.encodeString(v.String())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			enc.encodeBlob(v.Bytes())
			return
		}
		_ = enc.encodeList(v)
	case reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			if v.CanAddr() {
				b := v.Slice(0, v.Len()).Bytes()
				enc.encodeBlob(b)
				return
			}
			buf := make([]byte, v.Len())
			reflect.Copy(reflect.ValueOf(buf), v)
			enc.encodeBlob(buf)
			return
		}
		_ = enc.encodeList(v)
	default:
	}
}

// encode nil into one byte to buffer.
//
// +-----------+
// | 0000 1111 |   0x0f
// +-----------+
func (enc *Encoder) encodeNil() {
	enc.buf.WriteCode(Nil)
}

// 	Encode Boolean
// 	true
//
//	+-----------+
//	| 0000 0100 |   0x04
//	+-----------+
//	false
//
//  +-----------+
//  | 0000 0101 |   0x05
//  +-----------+
func (enc *Encoder) encodeBool(b bool) {
	if !b {
		enc.buf.WriteCode(False)
	} else {
		enc.buf.WriteCode(True)
	}
}

// Encode String
// String is also encoded into length + data like Blob.
//
// The type of String is 0x20, it also will be encoded into the last byte of the encoded bytes of length.
//
//              0x20 + 4 bits
// +...........+-----------+
// | 1xxx xxxx | 0010 xxxx |
// +...........+-----------+
func (enc *Encoder) encodeString(s string) {
	enc.encodeLen(len(s), String)
	enc.buf.WriteString(s)
}

// Encode Integer
// Except the last byte, the first bit of each byte will be 1.
// The remain 7 bits in these bytes and the remain 5 bits in the
// last byte will be used to store the value of the Integer.
//
//    7 bits                  5 bits
// +-----------+...........+-----------+
// | 1xxx xxxx | 1xxx xxxx | ...x xxxx |
// +-----------+...........+-----------+
func (enc *Encoder) encodeLen(n int, code Code) {
	for n > int(TagPackNumber) {
		enc.buf.WriteCode(NumSignBit | (Code(n) & NumMask))
		n >>= 7
	}
	enc.buf.WriteCode(code | Code(n))
}

// Encode Blob or []Byte
// Blob will be encoded into 2 parts. First part is the length of Blob, the second part is the binary data.
//
// +----------------+
// | length + data  |
// +----------------+
func (enc *Encoder) encodeBlob(b []byte) {
	enc.encodeLen(len(b), Blob)
	_, _ = enc.buf.Write(b)
}

// The Float type information will be encoded into the first byte,
// followed by bytes of the Float in the IEEE-754 format, in Big Endian.
//
// Double will be encoded into 9 bytes, Single will be 5 bytes.
//
//  0x06        8 bytes
// +-----------+===========+
// | 0000 0110 |    data   |  Double precision.
// +-----------+===========+
//
//  0x07        4 bytes
// +-----------+===========+
// | 0000 0111 |    data   |  Single precision.
// +-----------+===========+
func (enc *Encoder) encodeFloat32(f float32) {
	fb := math.Float32bits(f)
	enc.buf.WriteCode(Float)
	shift := byte(32)
	for shift > 0 {
		enc.buf.WriteCode(Code(fb & 0xff))
		fb >>= 8
		shift -= 8
	}
}

func (enc *Encoder) encodeFloat64(f float64) {
	fb := math.Float64bits(f)
	enc.buf.WriteCode(Double)
	shift := byte(64)
	for shift > 0 {
		enc.buf.WriteCode(Code(fb & 0xff))
		fb >>= 8
		shift -= 8
	}
}

// For encoding List and Dict, we define a Closure byte.
//
// +-----------+
// | 0000 0001 |   0x01, Closure
// +-----------+
// List type is encoded to one byte:
//
// +-----------+
// | 0000 0010 |   0x02, List
// +-----------+
// List type information will be encoded into the first byte, then following every element in List.
//
// The last byte is Closure.
//
// +-----------+
// | 0000 0010 |
// +-----------+----------------------------
// |          element 1
// +----------------------------------------
// |          element 2
// +----------------------------------------
// .    .    .
// .    .    .
// .    .    .
// +----------------------------------------
// |          element N
// +-----------+----------------------------
// | 0000 0001 | Closure
// +-----------+
func (enc *Encoder) encodeList(v reflect.Value) error {
	l := v.Len()
	enc.buf.WriteCode(List)
	for i := 0; i < l; i++ {
		if err := enc.EncodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	enc.buf.WriteCode(Closure)
	return nil
}
