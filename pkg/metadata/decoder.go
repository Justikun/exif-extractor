package metadata

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	"io"
)


func DecodeTagDataBytes(dataBytes []byte, dt DataType, count uint32, order binary.ByteOrder) (any, error) {
	//sliceReader := bytes.NewReader(dataBytes)
	//binaryReader := NewBinaryReader(sliceReader, order)

	//switch data
	return 0, fmt.Errorf("")

}

type BinaryReader struct {
	r 			io.ReadSeeker
	byteOrder 	binary.ByteOrder
}

func NewBinaryReader(r io.ReadSeeker, order binary.ByteOrder) *BinaryReader{
	return &BinaryReader{r: r, byteOrder: order}
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

func (br *BinaryReader) ReadSRational() (SRational, error) {
	var numerator 	int32
	var denominator int32

	err := binary.Read(br.r, br.byteOrder, &numerator)
	if err != nil { return SRational{}, err }
	err = binary.Read(br.r, br.byteOrder, &denominator)
	if err != nil { return SRational{}, err }
	return  SRational{Numerator: numerator, Denominator: denominator}, err
}

