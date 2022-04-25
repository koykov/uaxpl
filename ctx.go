package uaxpl

import (
	"github.com/koykov/bitset"
	"github.com/koykov/entry"
	"github.com/koykov/fastconv"
)

const (
	Unknown = "UNK"

	flagClientDetect = iota
	flagDeviceDetect
)

type Ctx struct {
	bitset.Bitset
	src []byte

	ct ClientType
	dt DeviceType

	be entry.Entry64
	ve entry.Entry64

	ctm ClientType
	dtm DeviceType
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

func (c *Ctx) GetBrowser() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.be > 0 {
		lo, hi := c.be.Decode()
		return fastconv.B2S(__cr_buf[lo:hi])
	}
	return Unknown
}

func (c *Ctx) GetBrowserVersionString() string {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	if c.ve > 0 {
		lo, hi := c.ve.Decode()
		return fastconv.B2S(c.src[lo:hi])
	}
	return Unknown
}

func (c *Ctx) GetClientType() ClientType {
	if !c.CheckBit(flagClientDetect) {
		c.parseClient()
	}
	return c.ct
}

func (c *Ctx) Reset() {
	c.reset()
	c.ctm = ClientTypeAll
	c.dtm = DeviceTypeAll
}

func (c *Ctx) reset() {
	c.Bitset.Reset()
	c.src = c.src[:0]
	c.be.Reset()
	c.ve.Reset()
}
