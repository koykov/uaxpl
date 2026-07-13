package uaxpl

import (
	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
)

type CacheEntry struct {
	bitset.Bitset

	ClientType      ClientType
	ClientName64    entry.Entry64
	ClientVersion64 entry.Entry64

	EngineName64    entry.Entry64
	EngineVersion64 entry.Entry64

	DeviceType  DeviceType
	BrandName64 entry.Entry64
	ModelName64 entry.Entry64

	OSName64    entry.Entry64
	OSVersion64 entry.Entry64

	Data []byte

	Hkey      uint64
	Timestamp int64
}

func (r *CacheEntry) FromCtx(ctx *Ctx) {
	r.Bitset = ctx.Bitset
	r.ClientType = ctx.clientType
	r.ClientName64 = ctx.clientName64
	r.ClientVersion64 = ctx.clientVersion64
	r.EngineName64 = ctx.engineName64
	r.EngineVersion64 = ctx.engineVersion64
	r.DeviceType = ctx.deviceType
	r.BrandName64 = ctx.brandName64
	r.ModelName64 = ctx.modelName64
	r.OSName64 = ctx.osName64
	r.OSVersion64 = ctx.osVersion64
	r.Data = bytealg.Copy(ctx.buf)
}
