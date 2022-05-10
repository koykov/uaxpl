package uaxpl

import (
	"bytes"

	"github.com/koykov/fastconv"
)

const (
	cpBrowserEngine = 0
	cpBrowser       = 1
	cpFeedReader    = 2
	cpLibrary       = 3
	cpMediaPlayer   = 4
	cpMobileApp     = 5
	cpPIM           = 6
)

func (c *Ctx) parseClient() bool {
	if c.ctm&ClientTypeBrowser != 0 {
		if c.evalClient(cpBrowser) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypeBrowser
			return true
		}
	}
	if c.ctm&ClientTypeFeedReader != 0 {
		if c.evalClient(cpFeedReader) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypeFeedReader
			return true
		}
	}
	if c.ctm&ClientTypeMobileApp != 0 {
		if c.evalClient(cpMobileApp) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypeMobileApp
			return true
		}
	}
	if c.ctm&ClientTypeMediaPlayer != 0 {
		if c.evalClient(cpMediaPlayer) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypeMediaPlayer
			return true
		}
	}
	if c.ctm&ClientTypePIM != 0 {
		if c.evalClient(cpPIM) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypePIM
			return true
		}
	}
	if c.ctm&ClientTypeLibrary != 0 {
		if c.evalClient(cpLibrary) {
			c.SetBit(flagClientDetect, true)
			c.ct = ClientTypeLibrary
			return true
		}
	}
	return false
}

func (c *Ctx) evalClient(idx int) bool {
	ir := __cr_idx[idx]
	irl := len(ir)
	_ = ir[irl-1]
	var x *cr
	for i := 0; i < irl; i++ {
		v := &ir[i]
		if v.re >= 0 {
			re := __cr_re[v.re]
			if re.Match(c.src) {
				x = v
				if x.vi != -1 {
					if m := re.FindSubmatchIndex(c.src); len(m) > int(x.vi) {
						lo1, hi1 := m[x.vi*2], m[x.vi*2+1]
						if lo1 != -1 && hi1 != -1 {
							lo, hi := uint32(m[x.vi*2]), uint32(m[x.vi*2+1])
							c.cve.Encode(lo, hi)
						}
					}
				}
				break
			}
		} else {
			lo, hi := v.si.Decode()
			si := __cr_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				break
			}
		}
	}

	if x != nil {
		c.cne = x.be

		// Engine detection.
		if idx == cpBrowser {
			if x.ed != 0 {
				c.ene = x.ed
			}
			if x.ef != -1 {
				lo, hi := c.cve.Decode()
				raw := c.src[lo:hi]
				c.ene = __cr_ef[x.ef](getMajor(fastconv.B2S(raw)))
			}
			if c.ene == 0 {
				ir = __cr_idx[cpBrowserEngine]
				irl = len(ir)
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
		}
		return true
	}
	return false
}
