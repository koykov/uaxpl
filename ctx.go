package uaxpl

type Ctx struct {
	src []byte
}

func NewCtx() *Ctx {
	ctx := Ctx{}
	return &ctx
}

func (c *Ctx) Reset() {
	c.src = c.src[:0]
}
