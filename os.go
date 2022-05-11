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
				if x.ni != -1 {
					if m := re.FindSubmatchIndex(c.src); len(m) > int(x.vi) {
						if x.ni != -1 {
							lo1, hi1 := m[x.ni*2], m[x.ni*2+1]
							c.os.Encode(uint32(lo1), uint32(hi1))
							c.SetBit(flagOSBufSrc, true)
						}
						// lo1, hi1 := m[x.vi*2], m[x.vi*2+1]
						// if lo1 != -1 && hi1 != -1 {
						// 	lo, hi := uint32(m[x.vi*2]), uint32(m[x.vi*2+1])
						// 	c.ove.Encode(lo, hi)
						// }
					}
				}
				// if x.vr != 0 {
				// 	lo, hi := x.vr.Decode()
				// 	rv := __or_ov[lo:hi]
				// 	rvl := len(rv)
				// 	_ = rv[rvl-1]
				// 	for j := lo; j < hi; j++ {
				// 		v1 := &__or_ov[j]
				// 		if v1.re >= 0 {
				// 			re1 := __or_re[v1.re]
				// 			if re1.Match(c.src) {
				// 				if v1.vi != -1 {
				// 					if m := re1.FindSubmatchIndex(c.src); len(m) > int(v1.vi) {
				// 						lo1, hi1 := m[v1.vi*2], m[v1.vi*2+1]
				// 						if lo1 != -1 && hi1 != -1 {
				// 							lo, hi := uint32(m[v1.vi*2]), uint32(m[v1.vi*2+1])
				// 							c.ove.Encode(lo, hi)
				// 							break
				// 						}
				// 					}
				// 				} else if v1.vs != 0 {
				// 					c.ove = v1.vs
				// 					c.SetBit(flagOSVerBufSrc, true)
				// 					break
				// 				}
				// 			}
				// 		} else if v1.si != 0 {
				// 			lo, hi := v1.si.Decode()
				// 			si := __or_buf[lo:hi]
				// 			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				// 				c.ove = v1.si
				// 				c.SetBit(flagOSVerBufSrc, true)
				// 				break
				// 			}
				// 		}
				// 	}
				// }
				break
			}
		} else if v.si != 0 {
			lo, hi := v.si.Decode()
			si := __or_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				break
			}
		}
	}

	if x != nil {
		if x.ni == -1 {
			c.os = x.ne
		}
		if x.vs != 0 {
			c.ove = x.vs
			c.SetBit(flagOSVerBufSrc, true)
		}

		c.SetBit(flagOSDetect, true)
		return true
	}

	return false
}