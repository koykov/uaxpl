package uavector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	cp ClientParser
	dp DeviceParser
}

// NewVector makes new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = uaHelper
	vec.cp, vec.dp = ClientParserAll, DeviceParserAll
	return vec
}

// Parse parses source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// ParseStr parses source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// ParseCopy makes a copy of source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// ParseCopyStr makes a copy of source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

func (vec *Vector) SetClientParser(mask ClientParser) *Vector {
	vec.cp = mask
	return vec
}

func (vec *Vector) SetDeviceParser(mask DeviceParser) *Vector {
	vec.dp = mask
	return vec
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.cp, vec.dp = ClientParserAll, DeviceParserAll
}
