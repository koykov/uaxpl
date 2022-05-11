package uaxpl

func NewCtxWithSrc(src []byte) *Ctx {
	ctx := NewCtx()
	return ctx.SetUserAgent(src)
}

func NewCtxWithSrcStr(src string) *Ctx {
	ctx := NewCtx()
	return ctx.SetUserAgentStr(src)
}

func AcquireWithSrc(src []byte) *Ctx {
	ctx := Acquire()
	return ctx.SetUserAgent(src)
}

func AcquireWithSrcStr(src string) *Ctx {
	ctx := Acquire()
	return ctx.SetUserAgentStr(src)
}

var _, _, _ = NewCtxWithSrc, NewCtxWithSrcStr, AcquireWithSrc
