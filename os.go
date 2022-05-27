package uaxpl

import (
	"bytes"

	"github.com/koykov/entry"
)

type osTuple struct {
	name64     entry.Entry64 // name
	nameSI     int8          // name match index
	matchRI    int32         // regex index
	match64    entry.Entry64 // substring
	versionSI  int8          // version match index
	version64  entry.Entry64 // static version
	versions64 entry.Entry64 // version ranges
}

type osVersionTuple struct {
	matchRI   int32         // regex index
	match64   entry.Entry64 // substring
	versionSI int8          // version match index
	version64 entry.Entry64 // static version
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
	var x *osTuple
	maybeMac := bytes.Index(c.src, osMaybeMacBytes[0]) != -1 ||
		bytes.Index(c.src, osMaybeMacBytes[1]) != -1 ||
		bytes.Index(c.src, osMaybeMacBytes[2]) != -1
	maybeNotAND := bytes.Index(c.src, osMaybeNotAND[0]) != -1
	for i := 0; i < rl; i++ {
		v := &r[i]
		if v.matchRI >= 0 {
			re := __or_re[v.matchRI]
			if rs := re.String(); (rs == osMaybeMacRE[0] || rs == osMaybeMacRE[1]) && maybeMac {
				continue
			}
			if rs := re.String(); (rs == osMaybeNotANDRE[0] || rs == osMaybeNotANDRE[1]) && maybeNotAND {
				continue
			}
			if re.Match(c.src) {
				x = v
				c.osName64 = x.name64
				if x.nameSI != -1 || x.versionSI != -1 {
					m := re.FindSubmatchIndex(c.src)
					if len(m) > int(x.nameSI) && x.nameSI != -1 {
						if lo1, hi1 := m[x.nameSI*2], m[x.nameSI*2+1]; lo1 != -1 && hi1 != -1 {
							c.osName64.Encode(uint32(lo1), uint32(hi1))
							c.SetBit(flagOSBufSrc, true)
						}
					}
					if len(m) > int(x.versionSI) && x.versionSI != -1 {
						if lo1, hi1 := m[x.versionSI*2], m[x.versionSI*2+1]; lo1 != -1 && hi1 != -1 {
							c.osVersion64.Encode(uint32(lo1), uint32(hi1))
						}
					}
				}
				if x.versions64 != 0 {
					c.osEvalVer(x.versions64)
				}
				break
			}
		} else if v.match64 != 0 {
			lo, hi := v.match64.Decode()
			si := __or_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				c.osName64 = x.name64
				break
			}
		}
	}

	if x != nil {
		if x.version64 != 0 {
			c.osVersion64 = x.version64
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
		if v.matchRI >= 0 {
			re := __or_re[v.matchRI]
			if re.Match(c.src) {
				if v.versionSI != -1 {
					if m := re.FindSubmatchIndex(c.src); len(m) > int(v.versionSI) {
						if lo1, hi1 := m[v.versionSI*2], m[v.versionSI*2+1]; lo1 != -1 && hi1 != -1 {
							c.osVersion64.Encode(uint32(lo1), uint32(hi1))
						}
					}
				} else {
					c.osVersion64 = v.version64
					c.SetBit(flagOSVerBufSrc, true)
				}
				break
			}
		} else if v.match64 != 0 {
			lo1, hi1 := v.match64.Decode()
			si := __or_buf[lo1:hi1]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				c.osVersion64 = v.version64
				c.SetBit(flagOSVerBufSrc, true)
				break
			}
		}
	}
}
