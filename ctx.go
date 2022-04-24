package uaxpl

type Ctx struct {
	src []byte

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
	c.src = append(c.src[:0], src...)
	return c
}

func (c *Ctx) SetUserAgentStr(src string) *Ctx {
	c.src = append(c.src[:0], src...)
	return c
}

func (c *Ctx) FilterClientType(mask ClientType) *Ctx {
	c.ctm = mask
	return c
}

func (c *Ctx) FilterDeviceType(mask DeviceType) *Ctx {
	c.dtm = mask
	return c
}

func (c *Ctx) Reset() {
	c.src = c.src[:0]
	c.ctm = ClientTypeAll
	c.dtm = DeviceTypeAll
}
