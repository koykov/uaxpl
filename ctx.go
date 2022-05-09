package uaxpl

import (
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
)

type Ctx struct {
	bitset.Bitset
	src []byte

	ctm ClientType
	dtm DeviceType

	hh uint64

	ct  ClientType
	cne entry.Entry64
	cve entry.Entry64
	cv  Version

	ene entry.Entry64
	eve entry.Entry64
	ev  Version

	dt DeviceType
	// ...
}

func NewCtx() *Ctx {
	ctx := Ctx{
		ctm: ClientTypeAll,
		dtm: DeviceTypeAll,
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
		c.hh = fnv.Hash64String(header)
	}
	return c
}

func (c *Ctx) SetUserAgentStr(src string) *Ctx {
	return c.SetUserAgent(fastconv.S2B(src))
}

func (c *Ctx) FilterClientType(mask ClientType) *Ctx {
	c.ctm = mask
	return c
}

func (c *Ctx) FilterDeviceType(mask DeviceType) *Ctx {
	c.dtm = mask
	return c
}

func (c *Ctx) GetUserAgent() string {
	return fastconv.B2S(c.src)
}

func (c *Ctx) GetClientType() ClientType {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	return c.ct
}

func (c *Ctx) GetBrowser() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	var (
		e   entry.Entry64
		buf []byte
	)
	if c.hh != 0 {
		if he, ok := __hr_idx[c.hh]; ok {
			e = he
			buf = __hr_buf
		}
	} else if c.cne > 0 {
		e = c.cne
		buf = __cr_buf
	}
	if e > 0 {
		lo, hi := e.Decode()
		return fastconv.B2S(buf[lo:hi])
	}
	return Unknown
}

func (c *Ctx) GetBrowserVersion() *Version {
	if !c.cv.p {
		raw := c.GetBrowserVersionString()
		_ = c.cv.Parse(raw)
	}
	return &c.cv
}

func (c *Ctx) GetBrowserVersionString() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.cve > 0 {
		lo, hi := c.cve.Decode()
		raw := c.src[lo:hi]
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
	if c.ene > 0 {
		lo, hi := c.ene.Decode()
		return fastconv.B2S(__cr_buf[lo:hi])
	}
	return Unknown
}

func (c *Ctx) GetEngineVersion() *Version {
	if !c.ev.p {
		raw := c.GetEngineVersionString()
		_ = c.ev.Parse(raw)
	}
	return &c.ev
}

func (c *Ctx) GetEngineVersionString() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.eve > 0 {
		lo, hi := c.eve.Decode()
		return fastconv.B2S(c.src[lo:hi])
	}
	return Unknown
}

func (c *Ctx) Reset() {
	c.reset()
	c.ctm = ClientTypeAll
	c.dtm = DeviceTypeAll
}

func (c *Ctx) reset() {
	c.Bitset.Reset()
	c.src = c.src[:0]

	c.hh = 0

	c.cne.Reset()
	c.cve.Reset()
	c.cv.Reset()

	c.ene.Reset()
	c.eve.Reset()
	c.ev.Reset()
}
