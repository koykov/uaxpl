package uaxpl

import (
	"bytes"
	"regexp"

	"github.com/koykov/byteconv"
)

var (
	forceEngineBlinkRE = regexp.MustCompile(`(?i)Chrome/.+ Safari/537.36`)
)

func (c *Ctx) evalEngine(cri *clientTuple) {
	if cri.engine64 != 0 {
		c.engineName64 = cri.engine64
	}
	if cri.engineFI != -1 {
		lo, hi := c.clientVersion64.Decode()
		raw := c.src[lo:hi]
		c.engineName64 = __cr_ef[cri.engineFI](getMajor(byteconv.B2S(raw)))
	}
	if c.engineName64 == 0 {
		ir := __cr_idx[cpBrowserEngine]
		irl := len(ir)
		_ = ir[irl-1]
		var e *clientTuple
		for i := 0; i < irl; i++ {
			v := &ir[i]
			if v.matchRI >= 0 {
				re := __cr_re[v.matchRI]
				if re.Match(c.src) {
					e = v
					break
				}
			} else {
				lo, hi := v.match64.Decode()
				si := __cr_buf[lo:hi]
				if len(si) > 0 && bytes.Index(c.src, si) != -1 {
					e = v
					break
				}
			}
		}
		if e != nil {
			c.engineName64 = e.browser64
		}
	}
	if c.engineName64 != 0 {
		if ri, ok := __cr_ev[c.engineName64]; ok {
			if m := __cr_evre[ri].FindSubmatchIndex(c.src); len(m) >= 4 {
				c.engineVersion64.Encode(uint32(m[2]), uint32(m[3]))
			}
		}
	}
	if c.engineName64 == 0 && forceEngineBlinkRE.Match(c.src) {
		c.SetBit(flagEngineForceBlink, true)
		return
	}
}

func (c *Ctx) evalEngineVersion() {}
