## EXIF Data Extractor

This project is in development as I learn more about EXIF, bytes, etc.

The goal is to create a few main functions.
- [0] Locate and extract EXIF APP1 APP0 marker.
  - Application marker APP1(0xFFE1) is used by JPEG
  - APP0(0xFFE0) is used by JFIF
- [0] Locate and Extract SOI marker (0xFFD8)
- [0] Extract APP1 Data (not to be confused with the APP1 Marker)
- [0] Extract EXIF Header (2 bytes)
- [0] Extract TIFF Header (2 bytes).
  - First 8 bytes of TIFF format images, are TIFF header.
  - 0x4949="II" is an "Intel" byte align.
  - 0x4d4d="MM" is a "Motorola" byte align.
- [0] Map tags to structure
  - e.g., 0x0132 is DateTime
  - I will need to find what the other tags are
- [0] Convert to JSON

As I learn more about EXIF and bytes I will try to update the information below.
By writing/explaining, it helps me retain new knowledge.
But maybe you will find an interest in it too!

### Byte Align
To help understand  little and big endian. We must first understand byte significance.

Take the number of 2125872. In hexidecimal, it's represented as 0x207030.
This number is 3 bytes (20, 70, 30).

Byte significance is determined by position. Similar to tenths, and hundreths place.
It is not determined by the numerical byte value.

THE MOST SIGNIFICANT is 20. Because it's the furthest byte to the left.
THE LEAST SIGNIFICANT is 30. Because it's the furstest byte to the left.

Different photo types (jpg, png, tiff, etc) store data with different endianness.
Endianness can be found with SOI marker (Start of Image).

#### Little-Endian Architecture - Intel
The least significant byte of a multi-byte value is stored at the lowest memory address.
The hexidecimal value of 0x12345678 (305,419,896) would be stored as 78 56 34 12.

#### Big-Endian Architecture - Motorola
The most significant byte is stored first, at the lowest memory address.
0x12345678 would be stored as 12 34 56 78.

- TIFF/EXIF: Little-endian or big-endian, based on header.
- GIF/BMP/WebP: Little-endian.
- JPG/PNG/JFIF/AVIF: Big-endian.
