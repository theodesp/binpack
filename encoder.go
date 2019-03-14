package binpack

import (
	"fmt"
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
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		enc.encodeInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		enc.encodeUInt(v)
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
		enc.encodeList(v)
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
		enc.encodeList(v)
	case reflect.Map:
		enc.encodeMap(v)
	default:
		enc.encodeString(fmt.Sprintf("binpack: Unsupported type %s", v.Type()))
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
func (enc *Encoder) encodeList(v reflect.Value) {
	l := v.Len()
	enc.buf.WriteCode(List)
	for i := 0; i < l; i++ {
		if err := enc.EncodeValue(v.Index(i)); err != nil {
			return
		}
	}
	enc.buf.WriteCode(Closure)
}

// Dict type:
//
// +-----------+
// | 0000 0011 |   0x03, Dict
// +-----------+
// Like List, the encoded data will begin with type information and end with
// Closure.
//
// The key and value of every Entry of the Dictionary will be encoded like
// following:
//
// +-----------+
// | 0000 0011 | Dict
// +-----------+----------------------------
// |            key 1
// +----------------------------------------
// |          value 1
// +----------------------------------------
// |            key 2
// +----------------------------------------
// |          value 2
// +----------------------------------------
// .    .    .
// .    .    .
// .    .    .
// +----------------------------------------
// |            key N
// +----------------------------------------
// |          value N
// +-----------+----------------------------
// | 0000 0001 | Closure
// +-----------+
func (enc *Encoder) encodeMap(v reflect.Value) {
	enc.buf.WriteCode(Dict)

	for _, key := range v.MapKeys() {
		if err := enc.EncodeValue(key); err != nil {
			return
		}
		if err := enc.EncodeValue(v.MapIndex(key)); err != nil {
			return
		}
	}
	enc.buf.WriteCode(Closure)
}

// Integer will be encoded into one or more bytes.
//
// The last byte is used to store the type and sign information of the Integer.
//
// The type and sign information is encode into the first 3 bits:
//
// positive
//
// +-----------+
// | 010x xxxx |  0x40
// +-----------+
// negative
//
// +-----------+
// | 011x xxxx |  0x60
// +-----------+
// Except the last byte, the first bit of each byte will be 1.
// The remain 7 bits in these bytes and the remain 5 bits in the
// last byte will be used to store the value of the Integer, For example:
//
//    7 bits                  5 bits
// +-----------+...........+-----------+
// | 1xxx xxxx | 1xxx xxxx | ...x xxxx |
// +-----------+...........+-----------+
func (enc *Encoder) encodeInt(v reflect.Value) {
	tag := Integer
	switch v.Kind() {
	case reflect.Int8:
		tag |= IntegerTypeByte
	case reflect.Int16:
		tag |= IntegerTypeShort
	case reflect.Int32:
		tag |= IntegerTypeInt
	case reflect.Int, reflect.Int64:
		tag |= IntegerTypeLong
	}

	val := v.Int()
	if val < 0 {
		val = -val
		tag |= IntegerNegative
	}

	for val > int64(TagPackInteger) || val>>3 > 0 {
		enc.buf.WriteCode(NumSignBit | (Code(val) & NumMask))
		val >>= 7
	}

	enc.buf.WriteCode(tag | Code(val))
}

func (enc *Encoder) encodeUInt(v reflect.Value) {
	tag := Integer
	switch v.Kind() {
	case reflect.Uint8:
		tag |= IntegerTypeByte
	case reflect.Uint16:
		tag |= IntegerTypeShort
	case reflect.Uint32:
		tag |= IntegerTypeInt
	case reflect.Uint, reflect.Uint64:
		tag |= IntegerTypeLong
	}

	val := v.Uint()
	for val > uint64(TagPackInteger) || val>>3 > 0 {
		enc.buf.WriteCode(NumSignBit | (Code(val) & NumMask))
		val >>= 7
	}
	enc.buf.WriteCode(tag | Code(val))
}
