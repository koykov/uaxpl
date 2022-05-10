package uaxpl

import (
	"bytes"

	"github.com/koykov/fastconv"
)

func (c *Ctx) evalEngine(cri *cr) {
	if cri.ed != 0 {
		c.ene = cri.ed
	}
	if cri.ef != -1 {
		lo, hi := c.cve.Decode()
		raw := c.src[lo:hi]
		c.ene = __cr_ef[cri.ef](getMajor(fastconv.B2S(raw)))
	}
	if c.ene == 0 {
		ir := __cr_idx[cpBrowserEngine]
		irl := len(ir)
		_ = ir[irl-1]
		var e *cr
		for i := 0; i < irl; i++ {
			v := &ir[i]
			if v.re >= 0 {
				re := __cr_re[v.re]
				if re.Match(c.src) {
					e = v
					break
				}
			} else {
				lo, hi := v.si.Decode()
				si := __cr_buf[lo:hi]
				if len(si) > 0 && bytes.Index(c.src, si) != -1 {
					e = v
					break
				}
			}
		}
		if e != nil {
			c.ene = e.be
		}
	}
	if c.ene != 0 {
		if ri, ok := __cr_ev[c.ene]; ok {
			if m := __cr_evre[ri].FindSubmatchIndex(c.src); len(m) >= 4 {
				c.eve.Encode(uint32(m[2]), uint32(m[3]))
			}
		}
	}
}

func (c *Ctx) evalEngineVersion() {}
