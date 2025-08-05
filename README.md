## EXIF Data Extractor

This project is in development as I learn more about EXIF, bytes, etc.

The goal is to create a few main functions.
- [] Locate and extract EXIF APP1 APP0 marker.
  - Application marker APP1(0xFFE1) is used by JPEG
  - APP0(0xFFE0) is used by JFIF
- [] Locate and Extract SOI marker (0xFFD8)
- [] Extract APP1 Data (not to be confused with the APP1 Marker)
- [] Extract EXIF Header (2 bytes)
- [] Extract TIFF Header (2 bytes).
  - First 8 bytes of TIFF format images, are TIFF header.
  - 0x4949="II" is an "Intel" byte align.
  - 0x4d4d="MM" is a "Motorola" byte align.
- [] Map tags to structure
  - e.g., 0x0132 is DateTime
  - I will need to find what the other tags are
- [] Convert to JSON 

As I learn more about EXIF and bytes I will try to update the information below.
By writing/explaining, it helps me retain new knowledge, but maybe you will find an interest in it too!

### Byte Align
To help understant little and big endian's. We must first understand byte significance.
Take the number of 2125872, in hexidecimal, it's represented as 0x207030. This number is 3 bytes (20, 70, 30).
Significance of bytes are determined by the position they sit in--similar to tenths, and hundreths place--and not the numerical value.

So our 3-byte number (0x207030 or split 20 70 30). 20 is THE MOST SIGNIFICANT byte because it's the furthest byte to the left, making 30 the least significant byte.
Anytime data is writen or read, we need to understand the endianness the data.Different photo types (jpg, png, tiff, etc) store data with different endianness. Endianness can be found with SOI, the Start of Image marker.

#### Little-Endian - Intel
The least significant byte of a multi-byte value is stored at the lowest memory address.
The hexidecimal value of 0x12345678 (represents this number 305,419,896) would be stored as 78 56 34 12

#### Big-Endian - Motorola
The most significant byte is stored first, at the lowest memory address.
0x12345678 would be stored as 12 34 56 78

