package uavector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
}

// NewVector makes new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = uaHelper
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
