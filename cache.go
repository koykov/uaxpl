package uaxpl

import (
	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
)

type Cacher[T any] interface {
	Set(key string, val T) error
	Get(string) (T, bool)
}

type cacheEntry struct {
	bitset.Bitset

	clientType      ClientType
	clientName64    entry.Entry64
	clientVersion64 entry.Entry64

	engineName64    entry.Entry64
	engineVersion64 entry.Entry64

	deviceType  DeviceType
	brandName64 entry.Entry64
	modelName64 entry.Entry64

	osName64    entry.Entry64
	osVersion64 entry.Entry64

	buf []byte

	hkey      uint64
	timestamp int64
}

func (r *cacheEntry) fromCtx(ctx *Ctx) {
	r.Bitset = ctx.Bitset
	r.clientType = ctx.clientType
	r.clientName64 = ctx.clientName64
	r.clientVersion64 = ctx.clientVersion64
	r.engineName64 = ctx.engineName64
	r.engineVersion64 = ctx.engineVersion64
	r.deviceType = ctx.deviceType
	r.brandName64 = ctx.brandName64
	r.modelName64 = ctx.modelName64
	r.osName64 = ctx.osName64
	r.osVersion64 = ctx.osVersion64
	r.buf = bytealg.Copy(ctx.buf)
}
