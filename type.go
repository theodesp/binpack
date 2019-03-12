package binpack

type Code byte

const (
	Closure Code = 0x01 // 0000 0001
	List    Code = 0x02 // 0000 0010
	Dict    Code = 0x03 // 0000 0011

	True  Code = 0x04 // 0000 0100
	False Code = 0x05 // 0000 0101

	Double Code = 0x06 // 0000 0110
	Float  Code = 0x07 // 0000 0111

	Nil    Code = 0x0f // 0000 1111
	Blob   Code = 0x10 // 0001 0000
	String Code = 0x20 // 0010 0000

	Integer         Code = 0x40 // 0100 0000
	IntegerNegative Code = 0x20 // 0010 0000

	IntegerTypeByte  Code = 0x01 << 3 // xxx0 1xxx
	IntegerTypeShort Code = 0x02 << 3 // xxx1 0xxx
	IntegerTypeInt   Code = 0x03 << 3 // xxx1 1xxx
	IntegerTypeLong  Code = 0x00 << 3 // xxx0 0xxx

	MaskIntegerSign      Code = 0x20 /* check if integer is negative */
	MaskTypeInteger      Code = 0x60 /* 0110 0000: integer or negative integer */
	MaskTypeStringOrBlob Code = 0x30 /* 00xx 0000: string or blob */
	MaskLastInteger      Code = 0x1f /* 000x xxxx the last 5 bits */
	MaskLastUintLen      Code = 0x0f /* 0000 xxxx the last 4 bits will be used to pack unit len */

	TagPackNumber  Code = 0x0f // 0001 xxxx
	TagPackInteger Code = 0x20 // 000x xxxx
	TagPackUintLen Code = 0x10 // 0000 xxxx

	NumSignBit Code = 0x80 // 1000 0000
	NumMask    Code = 0x7f // 0111 1111
)
