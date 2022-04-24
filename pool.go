package uaxpl

import "sync"

type Pool struct {
	p sync.Pool
}

var (
	// P is a default instance of the pool.
	// Just call uaxpl.Acquire() and uaxpl.Release().
	P Pool
	// Suppress go vet warnings.
	_, _ = Acquire, Release
)

// Get old context from the pool or create new one.
func (p *Pool) Get() *Ctx {
	v := p.p.Get()
	if v != nil {
		if ctx, ok := v.(*Ctx); ok {
			return ctx
		}
	}
	return NewCtx()
}

// Put context back to the pool.
func (p *Pool) Put(ctx *Ctx) {
	ctx.Reset()
	p.p.Put(ctx)
}

// Acquire gets context from default pool instance.
func Acquire() *Ctx {
	return P.Get()
}

// Release puts context back to default pool instance.
func Release(vec *Ctx) {
	P.Put(vec)
}
