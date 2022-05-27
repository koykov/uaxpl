package uaxpl

import (
	"bytes"
	"regexp"

	"github.com/koykov/entry"
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

func (c *Ctx) parseDevice() bool {
	if c.maskDeviceType&DeviceTypeCamera != 0 {
		if c.evalDevice(dpCamera) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeCamera
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeCarBrowser != 0 {
		if c.evalDevice(dpCarBrowser) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeCarBrowser
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeConsole != 0 {
		if c.evalDevice(dpConsole) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeConsole
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeMobile != 0 {
		if c.evalDevice(dpMobile) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeMobile
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeNotebook != 0 {
		if c.evalDevice(dpNotebook) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeNotebook
			return true
		}
	}
	if c.maskDeviceType&DeviceTypePortableMediaPlayer != 0 {
		if c.evalDevice(dpPortableMediaPlayer) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypePortableMediaPlayer
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeShellTV != 0 {
		if c.evalDevice(dpShellTV) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeShellTV
			return true
		}
	}
	if c.maskDeviceType&DeviceTypeTV != 0 {
		if c.evalDevice(dpTV) {
			c.SetBit(flagDeviceDetect, true)
			c.deviceType = DeviceTypeTV
			return true
		}
	}
	return false
}

func (c *Ctx) evalDevice(idx int) bool {
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
			c.deviceBufMNE(sm.model64)
		} else if x.models64 != 0 {
			lo, hi := x.models64.Decode()
			for i := lo; i < hi; i++ {
				m := &__dr_dm[i]
				if m.matchRI >= 0 {
					re := __dr_re[m.matchRI]
					if re.Match(c.src) {
						// todo replace RE placeholders
						c.deviceBufMNE1(m.model64, re)
						break
					}
				} else if m.match64 != 0 {
					lo1, hi1 := m.match64.Decode()
					si := __dr_buf[lo1:hi1]
					if len(si) > 0 && bytes.Index(c.src, si) != -1 {
						c.deviceBufMNE(m.model64)
						break
					}
				}
			}
		}
		return true
	}
	return false
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
	loop:
		p1 := raw[p+1]
		i := b2i(p1)
		r := c.src[m[i*2]:m[i*2+1]]
		c.buf = append(c.buf, r...)
		p = bytes.IndexByte(raw[p+1:], '$')
		if p != -1 {
			goto loop
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
