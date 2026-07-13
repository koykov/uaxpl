package uaxpl

import (
	"errors"
	"sync"
	"time"

	"github.com/koykov/hash/fnv"
)

const cacheTTL = int64(time.Hour)

type Cache interface {
	Set(key string, value *CacheEntry) error
	Get(key string) (*CacheEntry, error)
}

type cache struct {
	o   sync.Once
	mux sync.Mutex
	idx map[uint64]int
	buf []*CacheEntry
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

func (c *cache) Set(key string, row *CacheEntry) error {
	c.o.Do(c.init)

	row.Hkey = fnv.Hash64String(key)
	row.Timestamp = time.Now().UnixNano()

	c.mux.Lock()
	defer c.mux.Unlock()
	if idx, ok := c.idx[row.Hkey]; ok {
		c.buf[idx] = row
		return nil
	}
	c.buf = append(c.buf, row)
	c.idx[row.Hkey] = len(c.buf) - 1
	return nil
}

func (c *cache) Get(key string) (*CacheEntry, error) {
	hkey := fnv.Hash64String(key)

	c.o.Do(c.init)
	c.mux.Lock()
	defer c.mux.Unlock()
	if idx, ok := c.idx[hkey]; ok && idx >= 0 && idx < len(c.buf) {
		c.buf[idx].Timestamp = time.Now().UnixNano()
		return c.buf[idx], nil
	}
	return nil, errCache404
}

func (c *cache) clean() {
	now := time.Now().UnixNano()
	c.o.Do(c.init)
	c.mux.Lock()
	defer c.mux.Unlock()
	for i := 0; i < len(c.buf); i++ {
		if now-c.buf[i].Timestamp > cacheTTL {
			l := len(c.buf)
			old := c.buf[i].Hkey
			c.buf[i] = c.buf[l-1]
			c.buf = c.buf[:l-1]
			if i < len(c.buf) {
				// Edge case: has been deleted last item.
				c.idx[c.buf[i].Hkey] = i
			}
			delete(c.idx, old)
		}
	}
}

var errCache404 = errors.New("entry not found")
