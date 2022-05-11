package uaxpl

import (
	"bytes"

	"github.com/koykov/entry"
)

type or struct {
	ne entry.Entry64 // name
	ni int8          // name match index
	re int32         // regex index
	si entry.Entry64 // substring
	vi int8          // version match index
	vs entry.Entry64 // static version
	vr entry.Entry64 // version ranges
}

type ov struct {
	re int32         // regex index
	si entry.Entry64 // substring
	vi int8          // version match index
	vs entry.Entry64 // static version
}

var (
	osMaybeMacBytes = [][]byte{
		[]byte("Mac"),
		[]byte("(x86_64)"),
		[]byte("com.apple.Safari"),
	}
	osMaybeMacRE = []string{
		`(?i)^.*CFNetwork/.+ Darwin/(\d+[\.\d]+)`,
		`(?i)(?:Podcasts/(?:[\d\.]+)|Instacast(?:HD)?/(?:\d\.[\d\.abc]+)|Pocket Casts, iOS|\(iOS\)|iOS; Opera|Overcast|Castro|Podcat|iCatcher|RSSRadio/|MobileSafari/)`,
	}
	osMaybeNotAND = [][]byte{
		[]byte("like Android"),
	}
	osMaybeNotANDRE = []string{
		`(?i)(?:(?:Orca-)?Android|Adr|AOSP)[ /]?(?:[a-z]+ )?(\d+[\.\d]*)`,
		`(?i) Adr |Android|Silk-Accelerated=[a-z]{4,5}`,
	}
)

func (c *Ctx) parseOS() bool {
	r := __or_os
	rl := len(r)
	_ = r[rl-1]
	var x *or
	maybeMac := bytes.Index(c.src, osMaybeMacBytes[0]) != -1 ||
		bytes.Index(c.src, osMaybeMacBytes[1]) != -1 ||
		bytes.Index(c.src, osMaybeMacBytes[2]) != -1
	maybeNotAND := bytes.Index(c.src, osMaybeNotAND[0]) != -1
	for i := 0; i < rl; i++ {
		v := &r[i]
		if v.re >= 0 {
			re := __or_re[v.re]
			if rs := re.String(); (rs == osMaybeMacRE[0] || rs == osMaybeMacRE[1]) && maybeMac {
				continue
			}
			if rs := re.String(); (rs == osMaybeNotANDRE[0] || rs == osMaybeNotANDRE[1]) && maybeNotAND {
				continue
			}
			if re.Match(c.src) {
				x = v
				c.os = x.ne
				if x.ni != -1 || x.vi != -1 {
					m := re.FindSubmatchIndex(c.src)
					if len(m) > int(x.ni) && x.ni != -1 {
						if lo1, hi1 := m[x.ni*2], m[x.ni*2+1]; lo1 != -1 && hi1 != -1 {
							c.os.Encode(uint32(lo1), uint32(hi1))
							c.SetBit(flagOSBufSrc, true)
						}
					}
					if len(m) > int(x.vi) && x.vi != -1 {
						if lo1, hi1 := m[x.vi*2], m[x.vi*2+1]; lo1 != -1 && hi1 != -1 {
							c.ove.Encode(uint32(lo1), uint32(hi1))
						}
					}
				}
				if x.vr != 0 {
					c.osEvalVer(x.vr)
				}
				break
			}
		} else if v.si != 0 {
			lo, hi := v.si.Decode()
			si := __or_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				c.os = x.ne
				break
			}
		}
	}

	if x != nil {
		if x.vs != 0 {
			c.ove = x.vs
			c.SetBit(flagOSVerBufSrc, true)
		}

		c.SetBit(flagOSDetect, true)
		return true
	}

	return false
}

func (c *Ctx) osEvalVer(r entry.Entry64) {
	lo, hi := r.Decode()
	rv := __or_ov[lo:hi]
	rvl := len(rv)
	_ = rv[rvl-1]
	for i := 0; i < rvl; i++ {
		v := &rv[i]
		if v.re >= 0 {
			re := __or_re[v.re]
			if re.Match(c.src) {
				if v.vi != -1 {
					if m := re.FindSubmatchIndex(c.src); len(m) > int(v.vi) {
						if lo1, hi1 := m[v.vi*2], m[v.vi*2+1]; lo1 != -1 && hi1 != -1 {
							c.ove.Encode(uint32(lo1), uint32(hi1))
						}
					}
				} else {
					c.ove = v.vs
					c.SetBit(flagOSVerBufSrc, true)
				}
				break
			}
		} else if v.si != 0 {
			lo1, hi1 := v.si.Decode()
			si := __or_buf[lo1:hi1]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				c.ove = v.vs
				c.SetBit(flagOSVerBufSrc, true)
				break
			}
		}
	}
}
