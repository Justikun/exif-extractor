package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)


func DecodeTagData(dataBytes []byte, dt DataType, count uint32, order binary.ByteOrder) (any, error) {
	br := NewBinaryReader(bytes.NewReader(dataBytes), order)

	switch dt {
	case TypeByte:
		if count == 1 {
			val, err := br.ReadUint8()
			if err != nil { return nil, err }
			return val, nil
		}
		vals, err := br.ReadUint8()
		if err != nil { return nil, err }
		return vals, nil

	case TypeAscii:
		vals, err := br.ReadBytes(int(count))
		if err != nil {
	        return nil, fmt.Errorf("failed to read %d bytes for ASCII tag: %w", count, err)
		}
		s := string(vals)
		s = strings.TrimRight(s,"\x00 \t\r\n")
		return s, nil

	case TypeShort:
		vals := make([]uint16, count)
		for i := range vals {
			v, err := br.ReadUint16()
			if err != nil { return nil, err }
			vals[i] = v
		}
		return vals, nil

	case TypeLong:
		vals := make([]uint32, count)
		for i := range vals {
			v, err := br.ReadUint32()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	case TypeRational:
		 vals := make([]Rational, count)
	for i := range vals {
			v, err := br.ReadRational()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	case TypeSByte:
		vals := make([]int8, count)
		for i := range vals {
			v, err := br.ReadInt8()
			if err != nil { return nil, err }
			vals[i] = v
		}
		return vals, nil

	case TypeUndefined:
		// just raw bytes
		vals, err := br.ReadBytes(int(count))
		if err != nil { return nil, err }
		return vals, nil

	case TypeSShort:
		vals := make([]int16, count)
		for i:= range vals {
			v, err := br.ReadInt16()
			if err != nil { return nil, err }
			vals[i] = v
		}
		return vals, nil

	case TypeSLong:
		vals := make([]int32, count)
	for i := range vals {
			v, err := br.ReadInt32()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	case TypeSRational:
		vals := make([]Srational, count)
		for i := range vals {
			v, err := br.ReadSRational()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	case TypeFloat:
		vals := make([]float32, count)
		for i := range vals {
			v, err := br.ReadFloat32()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	case TypeDouble:
		vals := make([]float64, count)
		for i := range vals {
			v, err := br.ReadFloat64()
			if err != nil { return nil, err }
			vals[i] = v
		}
	return vals, nil

	default:
		return nil, fmt.Errorf("Unsupported datat type: %d", dt)
	}
}

type BinaryReader struct {
	r 			io.ReadSeeker
	byteOrder 	binary.ByteOrder
}

func NewBinaryReader(r io.ReadSeeker, order binary.ByteOrder) *BinaryReader{
	return &BinaryReader{ r: r, byteOrder: order }
}

func (br *BinaryReader) ChangeByteOrder(order binary.ByteOrder) {
	br.byteOrder = order
}

func (br *BinaryReader) Read(p []byte) (int, error) {
	return br.r.Read(p)
}

func (br *BinaryReader) ReadBytes(count int) ([]byte, error) {
	buf := make([]byte, count)
	n, err := br.r.Read(buf)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to read %v bytes: %w", count, err)
	}
	if n != count {
		return []byte{}, fmt.Errorf("Read %d bytes. Expected %d bytes", n, count)
	}
	return buf, err
}

func (br *BinaryReader) Seek(offset int64, whence int) (int64, error) {
	return br.r.Seek(offset, whence)
}

func (br *BinaryReader) SkipBytes(count int) error {
	_, err := br.r.Read(make([]byte, count))
	if err != nil {
		return fmt.Errorf("Failed to skip %v bytes", count)
	}
	return nil
}

// Specific binary readers
// Uint
func (br *BinaryReader) ReadUint8() (uint8, error) {
	var val uint8
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

func (br *BinaryReader) ReadUint16() (uint16, error) {
	var val uint16
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

func (br *BinaryReader) ReadUint32() (uint32, error) {
	var val uint32
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

// Int
func (br *BinaryReader) ReadInt8() (int8, error) {
	var val int8
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

func (br *BinaryReader) ReadInt16() (int16, error) {
	var val int16
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

func (br *BinaryReader) ReadInt32() (int32, error) {
	var val int32
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

// Float
func (br *BinaryReader) ReadFloat32() (float32, error) {
	var val float32
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

func (br *BinaryReader) ReadFloat64() (float64, error) {
	var val float64
	err := binary.Read(br.r, br.byteOrder, &val)
	return val, err
}

// Rational
func (br *BinaryReader) ReadRational() (Rational, error) {
	var numerator 	uint32
	var denominator uint32

	err := binary.Read(br.r, br.byteOrder, &numerator)
	if err != nil { return Rational{}, err }
	err = binary.Read(br.r, br.byteOrder, &denominator)
	if err != nil { return Rational{}, err }
	return Rational{Numerator: numerator, Denominator: denominator}, err
}

func (br *BinaryReader) ReadSRational() (Srational, error) {
	var numerator 	int32
	var denominator int32

	err := binary.Read(br.r, br.byteOrder, &numerator)
	if err != nil { return Srational{}, err }
	err = binary.Read(br.r, br.byteOrder, &denominator)
	if err != nil { return Srational{}, err }
	return  Srational{Numerator: numerator, Denominator: denominator}, err
}

