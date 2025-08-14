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

