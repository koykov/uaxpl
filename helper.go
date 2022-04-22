package uavector

import (
	"github.com/koykov/vector"
)

type UAHelper struct{}

var (
	uaHelper = &UAHelper{}
)

func (h *UAHelper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	return b
}
