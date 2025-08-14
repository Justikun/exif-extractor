[Resource](https://www.media.mit.edu/pia/Research/deepview/exif.html)
## JPEG
- Big Endian

Every jpeg file starts and end binary. Start and end binary have no data following the marker. It simply signifies the start and end.


### Flow of Markers
##### 1) Start of Image
- 0xFFD8 - SOI
##### 2) Meta Data Markers
**0xFFE1 - APP1 (Application Segment 1)**
Contains EXIF Data like
- camera settings
- date
- gps
- etc
**0xFFE0 - APP0 (Application Segment 0)**
Contains JFIF meta data
- version
- resolution
- thumbnail
- etc
**0xFFDB - DQT (Define Quantization Table)**
These tables specify how pixel data is quantized during JPG encoding
**0xFFC1 - SOF1 (Start of Frame 1)**
Defines the frame’s encoding parameters for extended sequential DCT, Huffman coding.
- image dimensions (width, height)
- number of color components
- sampling factors.
##### 3) Image Stream
- 0xFFDA - SOS
##### 4) End of Image
- 0xFFD9 - EOI

Please notice that "Data" contains Data size descriptor, if there is a Marker like this -> FF C1 00 0C

It means this Marker(0xFFC1) has 0x000C(equal 12)bytes of data. But the data size '12' includes "Data size" descriptor, it follows only 10 bytes of data after 0x000C.


Other binary markers are found in '0xFFXX', holding X data information about an image.
""

