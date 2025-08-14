package markers

const startHex = 0xFF

type Marker struct {
	Name	string
	First	byte
	Second 	byte
}

var SOIMarker = Marker  {
	Name: "SOI",
	First: startHex,
	Second: 0xD8,
}

var EOIMarker = Marker {
	Name: "EOI",
	First: startHex,
	Second: 0xD9,
}

var App0Marker = Marker {
	Name: "APP0",
	First: startHex,
	Second: 0xE0,
}

var App1Marker = Marker {
	Name: "APP1",
	First: startHex,
	Second: 0xE1,

}

var ImageStream = Marker {
	Name: "SOS",
	First: startHex,
	Second: 0xDA,
}
