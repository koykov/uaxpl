package uaxpl

import (
	"bytes"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
	"github.com/koykov/fastconv"
	"github.com/koykov/hash/fnv"
)

const (
	Unknown = "UNK"

	flagClientDetect = iota
	flagDeviceDetect
	flagOSDetect

	flagOSBufSrc
	flagOSVerBufSrc
	flagVerConstSrc

	flagDeviceForceDesktop
	flagEngineForceBlink
)

type Ctx struct {
	bitset.Bitset
	src, buf []byte

	maskClientType ClientType
	maskDeviceType DeviceType

	hintHash uint64

	clientType      ClientType
	clientName64    entry.Entry64
	clientVersion64 entry.Entry64
	clientVersion   Version

	engineName64    entry.Entry64
	engineVersion64 entry.Entry64
	engineVersion   Version

	deviceType  DeviceType
	brandName64 entry.Entry64
	modelName64 entry.Entry64

	osName64    entry.Entry64
	osVersion64 entry.Entry64
	osVersion   Version
}

func NewCtx() *Ctx {
	ctx := Ctx{
		maskClientType: ClientTypeAll,
		maskDeviceType: DeviceTypeAll,
	}
	return &ctx
}

func (c *Ctx) SetUserAgent(src []byte) *Ctx {
	c.reset()
	c.src = append(c.src[:0], src...)
	return c
}

func (c *Ctx) SetRequestedWith(header string) *Ctx {
	if len(header) > 0 {
		c.hintHash = fnv.Hash64String(header)
	}
	return c
}

func (c *Ctx) SetUserAgentStr(src string) *Ctx {
	return c.SetUserAgent(fastconv.S2B(src))
}

func (c *Ctx) FilterClientType(mask ClientType) *Ctx {
	c.maskClientType = mask
	return c
}

func (c *Ctx) FilterDeviceType(mask DeviceType) *Ctx {
	c.maskDeviceType = mask
	return c
}

func (c *Ctx) GetUserAgent() string {
	return fastconv.B2S(c.src)
}

func (c *Ctx) GetClientType() ClientType {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	return c.clientType
}

func (c *Ctx) GetBrowser() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	var (
		e   entry.Entry64
		buf []byte
	)
	if c.hintHash != 0 {
		if he, ok := __hr_idx[c.hintHash]; ok {
			e = he
			buf = __hr_buf
		}
	} else if c.clientName64 > 0 {
		e = c.clientName64
		buf = __cr_buf
	}
	if e > 0 {
		lo, hi := e.Decode()
		return fastconv.B2S(buf[lo:hi])
	}
	return Unknown
}

func (c *Ctx) GetBrowserVersion() *Version {
	if !c.clientVersion.p {
		raw := c.GetBrowserVersionString()
		_ = c.clientVersion.Parse(raw)
	}
	return &c.clientVersion
}

func (c *Ctx) GetBrowserVersionString() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.clientVersion64 > 0 {
		buf := c.src
		if c.CheckBit(flagVerConstSrc) {
			buf = __cr_buf
		}
		lo, hi := c.clientVersion64.Decode()
		raw := buf[lo:hi]
		if p := bytealg.IndexByteAtLR(raw, '/', 0); p != -1 {
			raw = raw[p+1:]
		}
		return fastconv.B2S(bytealg.TrimRight(raw, bDot))
	}
	return ""
}

func (c *Ctx) GetEngine() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.CheckBit(flagEngineForceBlink) {
		return "Blink"
	}
	if c.engineName64 > 0 {
		lo, hi := c.engineName64.Decode()
		return fastconv.B2S(__cr_buf[lo:hi])
	}
	return ""
}

func (c *Ctx) GetEngineVersion() *Version {
	if !c.engineVersion.p {
		raw := c.GetEngineVersionString()
		_ = c.engineVersion.Parse(raw)
	}
	return &c.engineVersion
}

func (c *Ctx) GetEngineVersionString() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.engineVersion64 > 0 {
		lo, hi := c.engineVersion64.Decode()
		return fastconv.B2S(c.src[lo:hi])
	}
	return ""
}

func (c *Ctx) GetDeviceType() DeviceType {
	if !c.CheckBit(flagDeviceDetect) {
		c.parseDevice()
	}
	if c.CheckBit(flagDeviceForceDesktop) {
		return DeviceTypeNotebook
	}
	return c.deviceType
}

func (c *Ctx) GetBrand() string {
	if !c.CheckBit(flagDeviceDetect) {
		c.parseDevice()
	}
	if c.brandName64 != 0 {
		lo, hi := c.brandName64.Decode()
		raw := __dr_buf[lo:hi]
		raw = bytealg.Trim(raw, bSpace)
		return fastconv.B2S(raw)
	}
	return Unknown
}

func (c *Ctx) GetModel() string {
	if !c.CheckBit(flagDeviceDetect) {
		c.parseDevice()
	}
	if c.modelName64 != 0 {
		lo, hi := c.modelName64.Decode()
		raw := c.buf[lo:hi]
		raw = bytealg.Trim(raw, bSpace)
		return fastconv.B2S(raw)
	}
	return ""
}

func (c *Ctx) GetOS() string {
	if !c.CheckBit(flagOSDetect) {
		c.parseOS()
	}
	if c.osName64 != 0 {
		buf := __or_buf
		if c.CheckBit(flagOSBufSrc) {
			buf = c.src
		}
		lo, hi := c.osName64.Decode()
		raw := buf[lo:hi]
		raw = bytealg.Trim(raw, bSpace)
		return fastconv.B2S(raw)
	}
	return Unknown
}

func (c *Ctx) GetOSVersion() *Version {
	if !c.osVersion.p {
		raw := c.GetOSVersionString()
		_ = c.osVersion.Parse(raw)
	}
	return &c.osVersion
}

func (c *Ctx) GetOSVersionString() string {
	if !c.CheckBit(flagOSDetect) {
		c.parseOS()
	}
	if c.osVersion64 > 0 {
		buf := c.src
		if c.CheckBit(flagOSVerBufSrc) {
			buf = __or_buf
		}
		lo, hi := c.osVersion64.Decode()
		raw := buf[lo:hi]
		if bytes.IndexByte(raw, '_') != -1 {
			off := len(c.buf)
			_ = raw[len(raw)-1]
			for i := 0; i < len(raw); i++ {
				if raw[i] == '_' {
					c.buf = append(c.buf, '.')
				} else {
					c.buf = append(c.buf, raw[i])
				}
			}
			raw = c.buf[off:]
		}
		return fastconv.B2S(raw)
	}
	return ""
}

func (c *Ctx) Reset() {
	c.reset()
	c.maskClientType = ClientTypeAll
	c.maskDeviceType = DeviceTypeAll
}

func (c *Ctx) reset() {
	c.Bitset.Reset()
	c.src = c.src[:0]
	c.buf = c.buf[:0]

	c.hintHash = 0

	c.clientType = 0
	c.clientName64.Reset()
	c.clientVersion64.Reset()
	c.clientVersion.Reset()

	c.engineName64.Reset()
	c.engineVersion64.Reset()
	c.engineVersion.Reset()

	c.deviceType = 0
	c.brandName64.Reset()
	c.modelName64.Reset()

	c.osName64.Reset()
	c.osVersion64.Reset()
	c.osVersion.Reset()
}

var (
	bSpace = []byte(" ")
)
