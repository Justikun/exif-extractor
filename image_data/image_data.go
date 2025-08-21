package image_data
import (
	"fmt"
	"encoding/binary"
)

type ImageData struct {
	ImagePath	string
	MetaData	MetaData
}

type MetaData struct {
	Tags		[]IFDTag
	ExifTags	[]IFDTag
	IntropTags	[]IFDTag
	GPSTags		[]IFDTag
}

type IFDTag struct {
	Id			uint16
	Name		string
	DataType	DataType
	Count		uint32
	Data		any
}

type Rational struct {
	Numerator   uint32
	Denominator uint32
}

type SRational struct {
	Numerator   int32
	Denominator int32
}

type IFDType string
const (
	IFDMAIN 		IFDType = "main"
	IFDEXIF			IFDType = "exif"
	IFDINTROP		IFDType = "introp"
	IFDGPS			IFDType = "gps"
)

func GetDataTypeFromBytes(b []byte, byteOrder binary.ByteOrder) (DataType, error) {
	dataValue := int(byteOrder.Uint16(b))
	switch DataType(dataValue) { // Cast the int to DataType for the switch
	case TypeByte:
		return TypeByte, nil
	case TypeAscii:
		return TypeAscii, nil
	case TypeShort:
		return TypeShort, nil
	case TypeLong:
		return TypeLong, nil
	case TypeRational:
		return TypeRational, nil
	case TypeSByte:
		return TypeSByte, nil
	case TypeUndefined:
		return TypeUndefined, nil
	case TypeSShort:
		return TypeSShort, nil
	case TypeSLong:
		return TypeSLong, nil
	case TypeSRational:
		return TypeSRational, nil
	case TypeFloat:
		return TypeFloat, nil
	case TypeDouble:
		return TypeDouble, nil
	default:
		return 0, fmt.Errorf("Failed to GetDataTypeBytes. Unknown data type value: %d", dataValue)
	}
}

func (dt DataType) GetDataTypeString() (string, error) {
	switch DataType(dt) { // Cast the int to DataType for the switch
	case TypeByte:
		return "BYTE", nil
	case TypeAscii:
		return "ASCII", nil
	case TypeShort:
		return "SHORT", nil
	case TypeLong:
		return "LONG", nil
	case TypeRational:
		return "RATIONAL", nil
	case TypeSByte:
		return "SBYTE", nil
	case TypeUndefined:
		return "UNDEFINED", nil
	case TypeSShort:
		return "SSHORT", nil
	case TypeSLong:
		return "SLONG", nil
	case TypeSRational:
		return "SRATIONAL", nil
	case TypeFloat:
		return "FLOAT", nil
	case TypeDouble:
		return "DOUBLE", nil
	default:
		return "", fmt.Errorf("Failed to GetDataTypeString. Unknown data type value: %d", dt)
	}
}

type DataType uint16

const (
	TypeByte      DataType = 1  // Unsigned 8-bit integer
	TypeAscii     DataType = 2  // 8-bit ASCII character
	TypeShort     DataType = 3  // Unsigned 16-bit integer
	TypeLong      DataType = 4  // Unsigned 32-bit integer
	TypeRational  DataType = 5  // Two LONGs (numerator, denominator)
	TypeSByte     DataType = 6  // Signed 8-bit integer
	TypeUndefined DataType = 7  // 8-bit untyped data
	TypeSShort    DataType = 8  // Signed 16-bit integer
	TypeSLong     DataType = 9  // Signed 32-bit integer
	TypeSRational DataType = 10 // Two SLONGs (signed numerator, denominator)
	TypeFloat     DataType = 11 // 32-bit IEEE floating point
	TypeDouble    DataType = 12 // 64-bit IEEE floating point
)

var TypeSize = map[DataType]uint32 {
	TypeByte:      1,
	TypeAscii:     1,
	TypeShort:     2,
	TypeLong:      4,
	TypeRational:  8,
	TypeSByte:     1,
	TypeUndefined: 1,
	TypeSShort:    2,
	TypeSLong:     4,
	TypeSRational: 8,
	TypeFloat:     4,
	TypeDouble:    8,
}

func (dt DataType) SizeInBytes() (int, error){
	switch dt {
	case TypeByte, TypeAscii, TypeSByte, TypeUndefined:
		return 1, nil
	case TypeShort, TypeSShort:
		return 2, nil
	case TypeLong, TypeSLong, TypeFloat:
		return 4, nil
	case TypeRational, TypeSRational, TypeDouble:
		return 8, nil
	default:
		return 0, fmt.Errorf("Unkown or invalid datatype for size calculation: %d", dt)
	}
}

func GetNameFromID(id uint16) (string, error) {
	if name, exists := ifdMainTagList[id]; exists {
		return name, nil
	}
	return "", fmt.Errorf("No Tag ID found")
}


var ifdMainTagList = map[uint16]string {
	// TIFF Baseline Tags (TIFF 6.0)
	0x00FE: "New Subfile Type",
	0x00FF: "Subfile Type",
	0x0100: "Image Width",
	0x0101: "Image Length",
	0x0102: "Bits Per Sample",
	0x0103: "Compression",
	0x0106: "Photometric Interpretation",
	0x0107: "Thresholding",
	0x0108: "Cell Width",
	0x0109: "Cell Length",
	0x010A: "Fill Order",
	0x010D: "Document Name",
	0x010E: "Image Description",
	0x010F: "Make",
	0x0110: "Model",
	0x0111: "Strip Offsets",
	0x0112: "Orientation",
	0x0115: "Samples Per Pixel",
	0x0116: "Rows Per Strip",
	0x0117: "Strip Byte Counts",
	0x0118: "Min Sample Value",
	0x0119: "Max Sample Value",
	0x011A: "X Resolution",
	0x011B: "Y Resolution",
	0x011C: "Planar Configuration",
	0x0128: "Resolution Unit",
	0x012D: "Transfer Function",
	0x0131: "Software",
	0x0132: "Date Time",
	0x013B: "Artist",
	0x013E: "White Point",
	0x013F: "Primary Chromaticities",
	// TIFF Extended Tags
	0x011D: "Page Name",
	0x011E: "X Position",
	0x011F: "Y Position",
	0x0120: "Free Offsets",
	0x0121: "Free Byte Counts",
	0x0122: "Gray Response Unit",
	0x0123: "Gray Response Curve",
	0x0124: "T4 Options",
	0x0125: "T6 Options",
	0x0129: "Page Number",
	0x013C: "Host Computer",
	0x013D: "Predictor",
	0x0140: "Color Map",
	0x0141: "Halftone Hints",
	0x0142: "Tile Width",
	0x0143: "Tile Length",
	0x0144: "Tile Offsets",
	0x0145: "Tile Byte Counts",
	0x014C: "Ink Set",
	0x014D: "Ink Names",
	0x014E: "Number Of Inks",
	0x0150: "Dot Range",
	0x0151: "Target Printer",
	0x0152: "Extra Samples",
	0x0153: "Sample Format",
	0x0154: "S Min Sample Value",
	0x0155: "S Max Sample Value",
	0x0156: "Transfer Range",
	0x0200: "JPEG Proc",
	0x0201: "JPEG Interchange Format",
	0x0202: "JPEG Interchange Format Length",
	0x0203: "JPEG Restart Interval",
	0x0205: "JPEG Lossless Predictors",
	0x0206: "JPEG Point Transforms",
	0x0207: "JPEG Q Tables",
	0x0208: "JPEG DC Tables",
	0x0209: "JPEG AC Tables",
	0x0211: "Y Cb Cr Coefficients",
	0x0212: "Y Cb Cr Sub Sampling",
	0x0213: "Y Cb Cr Positioning",
	0x0214: "Reference Black White",
	// Other Common Tags
	0x8769: "Exif IFD Pointer",
	0x8825: "GPS Info IFD Pointer",
	0x02BC: "Application Notes", // XMP metadata
	0x83BB: "IPTC Data",
	0x8773: "ICC Profile",
}

var ifdExifTagList = map[uint16]string {
	0x8298: "Copyright",
	0x829A: "Exposure Time",
	0x829D: "F Number",
	0x8822: "Exposure Program",
	0x8824: "Spectral Sensitivity",
	0x8827: "ISO Speed Ratings",
	0x8828: "Opto-Electric Conversion Factor",
	0x882A: "Time Zone Offset",
	0x882B: "Self Timer Mode",
	0x9000: "Exif Version",
	0x9003: "Date Time Original",
	0x9004: "Create Date",
	0x9101: "Components Configuration",
	0x9102: "Compressed Bits Per Pixel",
	0x9201: "Shutter Speed Value",
	0x9202: "Aperture Value",
	0x9203: "Brightness Value",
	0x9204: "Exposure Compensation",
	0x9205: "Max Aperture Value",
	0x9206: "Subject Distance",
	0x9207: "Metering Mode",
	0x9208: "Light Source",
	0x9209: "Flash",
	0x920A: "Focal Length",
	0x9214: "Subject Area",
	0x927C: "Maker Note",
	0x9286: "User Comment",
	0x9290: "Sub Sec Time",
	0x9291: "Sub Sec Time Original",
	0x9292: "Sub Sec Time Digitized",
	0xA000: "Flashpix Version",
	0xA001: "Color Space",
	0xA002: "Pixel X Dimension",
	0xA003: "Pixel Y Dimension",
	0xA004: "Related Sound File",
	0xA005: "Interop Offset",
	0xA20B: "Flash Energy",
	0xA20C: "Spatial Frequency Response",
	0xA20E: "Focal Plane X Resolution",
	0xA20F: "Focal Plane Y Resolution",
	0xA210: "Focal Plane Resolution Unit",
	0xA214: "Subject Location",
	0xA215: "Exposure Index",
	0xA217: "Sensing Method",
	0xA300: "File Source",
	0xA301: "Scene Type",
	0xA302: "CFA Pattern",
	0xA401: "Custom Rendered",
	0xA402: "Exposure Mode",
	0xA403: "White Balance",
	0xA404: "Digital Zoom Ratio",
	0xA405: "Focal Length In 35mm Film",
	0xA406: "Scene Capture Type",
	0xA407: "Gain Control",
	0xA408: "Contrast",
	0xA409: "Saturation",
	0xA40A: "Sharpness",
	0xA40C: "Subject Distance Range",
	0xA420: "Image Unique ID",
}

var ifdIntropTagList = map[uint16]string {
	0x0001: "Interop Index",
	0x0002: "Interop Version",
	0x1000: "Related Image File Format",
	0x1001: "Related Image Width",
	0x1002: "Related Image Length",
}

var ifdGPSTagList = map[uint16]string {
	0x0000: "GPS Version ID",
	0x0001: "GPS Latitude Ref",
	0x0002: "GPS Latitude",
	0x0003: "GPS Longitude Ref",
	0x0004: "GPS Longitude",
	0x0005: "GPS Altitude Ref",
	0x0006: "GPS Altitude",
	0x0007: "GPS Time Stamp",
	0x0008: "GPS Satellites",
	0x0009: "GPS Status",
	0x000A: "GPS Measure Mode",
	0x000B: "GPS DOP",
	0x000C: "GPS Speed Ref",
	0x000D: "GPS Speed",
	0x000E: "GPS Track Ref",
	0x000F: "GPS Track",
	0x0010: "GPS Img Direction Ref",
	0x0011: "GPS Img Direction",
	0x0012: "GPS Map Datum",
	0x0013: "GPS Dest Latitude Ref",
	0x0014: "GPS Dest Latitude",
	0x0015: "GPS Dest Longitude Ref",
	0x0016: "GPS Dest Longitude",
	0x0017: "GPS Dest Bearing Ref",
	0x0018: "GPS Dest Bearing",
	0x0019: "GPS Dest Distance Ref",
	0x001A: "GPS Dest Distance",
	0x001B: "GPS Processing Method",
	0x001C: "GPS Area Information",
	0x001D: "GPS Date Stamp",
	0x001E: "GPS Differential",
}
