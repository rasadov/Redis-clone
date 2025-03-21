package rdb

import "io"

// buffer is a trivial growable byte buffer.
type buffer struct {
	bytes []byte
}

func (b *buffer) writeByte(x byte) {
	b.bytes = append(b.bytes, x)
}
func (b *buffer) write(p []byte) {
	b.bytes = append(b.bytes, p...)
}

// encodeSize writes a size < 64 as a single byte: 2 bits of 00 + 6 bits of the size.
func encodeSize(b *buffer, n uint64) {
	if n > 63 {
		panic("encodeSize only handles n <= 63 in this minimal example")
	}
	b.writeByte(byte(n & 0x3F))
}

// encodeString writes a size-encoded string (length ≤ 63) plus the raw bytes.
func encodeString(b *buffer, s string) {
	if len(s) > 63 {
		panic("encodeString only handles string length ≤ 63 in this minimal example")
	}
	// encode size
	encodeSize(b, uint64(len(s)))
	// then raw chars
	b.write([]byte(s))
}

// decodeSize is the inverse of encodeSize (for n ≤ 63).
// Returns (value, newIndex, error).
func decodeSize(Data []byte, start int) (uint64, int, error) {
	if start >= len(Data) {
		return 0, start, io.ErrUnexpectedEOF
	}
	b := Data[start]
	// top two bits?
	// in this mini-implementation, we only handle 0b00 with 6 bits of Data
	top2 := b >> 6
	if top2 != 0b00 {
		// not handled in this snippet
		return 0, start, nil
	}
	val := uint64(b & 0x3F)
	return val, start + 1, nil
}

// decodeString is the inverse of encodeString (assuming length ≤ 63).
func decodeString(Data []byte, start int) (string, int, error) {
	length, newPos, err := decodeSize(Data, start)
	if err != nil {
		return "", newPos, err
	}
	end := newPos + int(length)
	if end > len(Data) {
		return "", newPos, io.ErrUnexpectedEOF
	}
	return string(Data[newPos:end]), end, nil
}
