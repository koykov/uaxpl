package uaxpl

import (
	"bytes"
	"regexp"

	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

const (
	dpCamera              = 0
	dpCarBrowser          = 1
	dpConsole             = 2
	dpMobile              = 3
	dpNotebook            = 4
	dpPortableMediaPlayer = 5
	dpShellTV             = 6
	dpTV                  = 7
)

var (
	desktopOS = map[string]struct{}{
		"AmigaOS":     {},
		"IBM":         {},
		"GNU/Linux":   {},
		"Ubuntu":      {},
		"Fedora":      {},
		"Mac":         {},
		"Unix":        {},
		"OpenBSD":     {},
		"FreeBSD":     {},
		"Windows":     {},
		"BeOS":        {},
		"Chrome OS":   {},
		"Chromium OS": {},
		"OS ROSA":     {},
	}

	reForceTablet = regexp.MustCompile(`(?i)Tablet`)
	reForceMobile = regexp.MustCompile(`(?i)Mobile`)
)

func (c *Ctx) parseDevice() bool {
	if c.maskDeviceType&DeviceTypeCamera != 0 {
		if typ, ok := c.evalDevice(dpCamera, DeviceTypeCamera); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeCarBrowser != 0 {
		if typ, ok := c.evalDevice(dpCarBrowser, DeviceTypeCarBrowser); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeConsole != 0 {
		if typ, ok := c.evalDevice(dpConsole, DeviceTypeConsole); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypePortableMediaPlayer != 0 {
		if typ, ok := c.evalDevice(dpPortableMediaPlayer, DeviceTypeNotebook); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeShellTV != 0 {
		if typ, ok := c.evalDevice(dpShellTV, DeviceTypeShellTV); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeTV != 0 {
		if typ, ok := c.evalDevice(dpTV, DeviceTypeTV); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeMobile != 0 {
		if typ, ok := c.evalDevice(dpMobile, DeviceTypeMobile); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeNotebook != 0 {
		if typ, ok := c.evalDevice(dpNotebook, DeviceTypeNotebook); ok {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = typ
			return true
		}
	}

	if reForceTablet.Match(c.src) {
		c.deviceType = DeviceTypeTablet
		return true
	}

	if reForceMobile.Match(c.src) {
		c.deviceType = DeviceTypeMobile
		return true
	}

	if c.brandName64 == 0 {
		c.evalVendor()
	}

	return false
}

func (c *Ctx) evalDevice(idx int, defType DeviceType) (typ DeviceType, ok bool) {
	typ = defType

	defer func() {
		if ok || idx != dpNotebook {
			return
		}
		os := c.GetOS()
		if _, ok1 := desktopOS[os]; ok1 {
			c.SetBit(flagDeviceForceDesktop, true)
		}
	}()

	ir := __dr_idx[idx]
	irl := len(ir)
	_ = ir[irl-1]
	var x *deviceTuple
	for i := 0; i < irl; i++ {
		v := &ir[i]
		if v.matchRI >= 0 {
			re := __dr_re[v.matchRI]
			if re.Match(c.src) {
				x = v
				break
			}
		} else {
			lo, hi := v.match64.Decode()
			si := __dr_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				break
			}
		}
	}
	if x != nil {
		c.brandName64 = x.brand64

		if x.modelSI != -1 {
			sm := &__dr_dm[x.modelSI]
			if x.matchRI != -1 {
				c.deviceBufMNE1(sm.model64, __dr_re[x.matchRI])
			} else {
				c.deviceBufMNE(sm.model64)
			}
			typ = c.deviceEvalType(sm.type64, x.type64, defType)
		} else if x.models64 != 0 {
			lo, hi := x.models64.Decode()
			for i := lo; i < hi; i++ {
				m := &__dr_dm[i]
				if m.matchRI >= 0 {
					re := __dr_re[m.matchRI]
					if re.Match(c.src) {
						c.deviceBufMNE1(m.model64, re)
						typ = c.deviceEvalType(m.type64, x.type64, defType)
						break
					}
				} else if m.match64 != 0 {
					lo1, hi1 := m.match64.Decode()
					si := __dr_buf[lo1:hi1]
					if len(si) > 0 && bytes.Index(c.src, si) != -1 {
						c.deviceBufMNE(m.model64)
						typ = c.deviceEvalType(m.type64, x.type64, defType)
						break
					}
				}
			}
		} else {
			typ = c.deviceEvalType(x.type64, 0, defType)
		}
		ok = true
	}
	return
}

func (c *Ctx) deviceBufMNE(e entry.Entry64) {
	lo, hi := e.Decode()
	raw := __dr_buf[lo:hi]
	lo1 := uint32(len(c.buf))
	c.buf = append(c.buf, raw...)
	hi1 := uint32(len(c.buf))
	c.modelName64.Encode(lo1, hi1)
}

func (c *Ctx) deviceBufMNE1(e entry.Entry64, re *regexp.Regexp) {
	lo, hi := e.Decode()
	raw := __dr_buf[lo:hi]
	p := bytes.IndexByte(raw, '$')
	if p != -1 {
		lo1 := uint32(len(c.buf))
		c.buf = append(c.buf, raw[:p]...)
		m := re.FindSubmatchIndex(c.src)
		for {
			p1 := raw[p+1]
			i := b2i(p1)
			if len(m) <= i*2+1 {
				break
			}
			lo2, hi2 := m[i*2], m[i*2+1]
			if lo2 < 0 || hi2 < 0 {
				break
			}
			r := c.src[m[i*2]:m[i*2+1]]
			c.buf = append(c.buf, r...)
			pp := p + 1
			p = bytealg.IndexByteAtLR(raw, '$', pp)
			if p-(pp+1) > 0 {
				c.buf = append(c.buf, raw[pp+1:p]...)
			}
			if p != -1 {
				continue
			}
			c.buf = append(c.buf, raw[pp+1:]...)
			break
		}
		hi1 := uint32(len(c.buf))
		c.modelName64.Encode(lo1, hi1)
		return
	}
	lo1 := uint32(len(c.buf))
	c.buf = append(c.buf, raw...)
	hi1 := uint32(len(c.buf))
	c.modelName64.Encode(lo1, hi1)
}

func (c *Ctx) deviceEvalType(typ, typ1 entry.Entry64, defType DeviceType) DeviceType {
	lo, hi := typ.Decode()
	raw := fastconv.B2S(__dr_buf[lo:hi])
	if len(raw) == 0 {
		lo, hi = typ1.Decode()
	}
	raw = fastconv.B2S(__dr_buf[lo:hi])
	if len(raw) == 0 {
		return defType
	}
	switch raw {
	case "smartphone":
		return DeviceTypeMobile
	case "phablet":
		return DeviceTypePhablet
	case "tablet":
		return DeviceTypeTablet
	case "tv":
		return DeviceTypeTV
	case "console":
		return DeviceTypeConsole
	case "portable media player":
		return DeviceTypePortableMediaPlayer
	default:
		return defType
	}
}

func (c *Ctx) evalVendor() {
	for i := 0; i < len(__vr_idx); i++ {
		brand := __vr_idx[i]
		lo, hi := brand.range64.Decode()
		for j := lo; j < hi; j++ {
			re := __vr_re[j]
			if re.Match(c.src) {
				c.brandName64 = brand.brand64
				c.SetBit(flagVendorFragment, true)
			}
		}
	}
}
