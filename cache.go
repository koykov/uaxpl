package uaxpl

import (
	"sync"
	"time"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
	"github.com/koykov/hash/fnv"
)

const (
	cacheTTL = int64(time.Hour)
)

type cache struct {
	o   sync.Once
	mux sync.Mutex
	idx map[uint64]int
	buf []cacheRow
}

type cacheRow struct {
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

func (c *cache) init() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.idx = make(map[uint64]int)
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.clean()
			}
		}
	}()
}

func (c *cache) set(key string, row cacheRow) {
	c.o.Do(c.init)

	row.hkey = fnv.Hash64String(key)
	row.timestamp = time.Now().UnixNano()

	c.mux.Lock()
	defer c.mux.Unlock()
	if idx, ok := c.idx[row.hkey]; ok {
		c.buf[idx] = row
		return
	}
	c.buf = append(c.buf, row)
	c.idx[row.hkey] = len(c.buf) - 1
}

func (c *cache) get(key string) (*cacheRow, bool) {
	hkey := fnv.Hash64String(key)

	c.o.Do(c.init)
	c.mux.Lock()
	defer c.mux.Unlock()
	if idx, ok := c.idx[hkey]; ok && idx >= 0 && idx < len(c.buf) {
		c.buf[idx].timestamp = time.Now().UnixNano()
		return &c.buf[idx], true
	}
	return nil, false
}

func (c *cache) clean() {
	now := time.Now().UnixNano()
	c.o.Do(c.init)
	c.mux.Lock()
	defer c.mux.Unlock()
	for i := 0; i < len(c.buf); i++ {
		if now-c.buf[i].timestamp > cacheTTL {
			old := c.buf[i].hkey
			c.buf[i] = c.buf[len(c.buf)-1]
			c.buf = c.buf[:len(c.buf)-1]
			c.idx[c.buf[i].hkey] = i
			delete(c.idx, old)
		}
	}
}

func (r *cacheRow) fromCtx(ctx *Ctx) {
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
