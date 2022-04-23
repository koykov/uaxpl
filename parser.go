package uaxpl

import (
	"github.com/koykov/bytealg"
)

var (
	// Byte constants.
	bSpace = []byte(" ")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bSpace)

	return
}
