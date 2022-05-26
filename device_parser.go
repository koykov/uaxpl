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
	if c.dtm&DeviceTypeCamera != 0 {
		if c.evalDevice(dpCamera) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeCamera
			return true
		}
	}
	if c.dtm&DeviceTypeCarBrowser != 0 {
		if c.evalDevice(dpCarBrowser) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeCarBrowser
			return true
		}
	}
	if c.dtm&DeviceTypeConsole != 0 {
		if c.evalDevice(dpConsole) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeConsole
			return true
		}
	}
	if c.dtm&DeviceTypeMobile != 0 {
		if c.evalDevice(dpMobile) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeMobile
			return true
		}
	}
	if c.dtm&DeviceTypeNotebook != 0 {
		if c.evalDevice(dpNotebook) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeNotebook
			return true
		}
	}
	if c.dtm&DeviceTypePortableMediaPlayer != 0 {
		if c.evalDevice(dpPortableMediaPlayer) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypePortableMediaPlayer
			return true
		}
	}
	if c.dtm&DeviceTypeShellTV != 0 {
		if c.evalDevice(dpShellTV) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeShellTV
			return true
		}
	}
	if c.dtm&DeviceTypeTV != 0 {
		if c.evalDevice(dpTV) {
			c.SetBit(flagDeviceDetect, true)
			c.dt = DeviceTypeTV
			return true
		}
	}
	return false
}

func (c *Ctx) evalDevice(idx int) bool {
	ir := __dr_idx[idx]
	irl := len(ir)
	_ = ir[irl-1]
	var x *dr
	for i := 0; i < irl; i++ {
		v := &ir[i]
		if v.re >= 0 {
			re := __dr_re[v.re]
			if re.Match(c.src) {
				x = v
				break
			}
		} else {
			lo, hi := v.si.Decode()
			si := __dr_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				break
			}
		}
	}
	if x != nil {
		c.bne = x.ne

		if x.sm != -1 {
			sm := &__dr_dm[x.sm]
			c.deviceBufMNE(sm.ne)
		} else if x.me != 0 {
			lo, hi := x.me.Decode()
			for i := lo; i < hi; i++ {
				m := &__dr_dm[i]
				if m.re >= 0 {
					re := __dr_re[m.re]
					if re.Match(c.src) {
						// todo replace RE placeholders
						c.deviceBufMNE1(m.ne, re)
						break
					}
				} else if m.si != 0 {
					lo1, hi1 := m.si.Decode()
					si := __dr_buf[lo1:hi1]
					if len(si) > 0 && bytes.Index(c.src, si) != -1 {
						c.deviceBufMNE(m.ne)
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
	c.mne.Encode(lo1, hi1)
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
		c.mne.Encode(lo1, hi1)
		return
	}
	lo1 := uint32(len(c.buf))
	c.buf = append(c.buf, raw...)
	hi1 := uint32(len(c.buf))
	c.mne.Encode(lo1, hi1)
}
