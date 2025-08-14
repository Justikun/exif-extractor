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

