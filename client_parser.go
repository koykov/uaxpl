package uaxpl

import (
	"bytes"
)

const (
	cpBrowserEngine = 0
	cpBrowser       = 1
	cpFeedReader    = 2
	cpLibrary       = 3
	cpMediaPlayer   = 4
	cpMobileApp     = 5
	cpPIM           = 6
)

func (c *Ctx) parseClient() (ok bool) {
	// Check cached result.
	if row, ok1 := cache_.get(c.GetUserAgent()); ok1 && row.CheckBit(flagClientDetect) {
		c.fromCache(row)
		ok = true
		return
	}
	defer func() {
		c.SetBit(flagClientDetect, true)
		if ok {
			// Put result to the cache.
			var row cacheRow
			row.fromCtx(c)
			cache_.set(c.GetUserAgent(), row)
		}
	}()

	if c.maskClientType&ClientTypeFeedReader != 0 {
		if c.evalClient(cpFeedReader) {
			c.clientType = ClientTypeFeedReader
			ok = true
			return
		}
	}
	if c.maskClientType&ClientTypeMobileApp != 0 {
		if c.evalClient(cpMobileApp) {
			c.clientType = ClientTypeMobileApp
			ok = true
			return
		}
	}
	if c.maskClientType&ClientTypeMediaPlayer != 0 {
		if c.evalClient(cpMediaPlayer) {
			c.clientType = ClientTypeMediaPlayer
			ok = true
			return
		}
	}
	if c.maskClientType&ClientTypePIM != 0 {
		if c.evalClient(cpPIM) {
			c.clientType = ClientTypePIM
			ok = true
			return
		}
	}
	if c.maskClientType&ClientTypeLibrary != 0 {
		if c.evalClient(cpLibrary) {
			c.clientType = ClientTypeLibrary
			ok = true
			return
		}
	}
	if c.maskClientType&ClientTypeBrowser != 0 {
		if c.evalClient(cpBrowser) {
			c.clientType = ClientTypeBrowser
			ok = true
			return
		}
	}
	return
}

func (c *Ctx) evalClient(idx int) bool {
	ir := __cr_idx[idx]
	irl := len(ir)
	_ = ir[irl-1]
	var x *clientTuple
	for i := 0; i < irl; i++ {
		v := &ir[i]
		if v.matchRI >= 0 {
			re := __cr_re[v.matchRI]
			if re.Match(c.src) {
				x = v
				if x.version64 != 0 {
					c.clientVersion64 = x.version64
					c.SetBit(flagVerConstSrc, true)
				} else if x.browserVI != -1 {
					if m := re.FindSubmatchIndex(c.src); len(m) > int(x.browserVI) {
						lo1, hi1 := m[x.browserVI*2], m[x.browserVI*2+1]
						if lo1 != -1 && hi1 != -1 {
							lo, hi := uint32(m[x.browserVI*2]), uint32(m[x.browserVI*2+1])
							c.clientVersion64.Encode(lo, hi)
						}
					}
				}
				break
			}
		} else {
			lo, hi := v.match64.Decode()
			si := __cr_buf[lo:hi]
			if len(si) > 0 && bytes.Index(c.src, si) != -1 {
				x = v
				break
			}
		}
	}

	if x != nil {
		c.clientName64 = x.browser64

		// Engine detection.
		if idx == cpBrowser {
			c.evalEngine(x)
		}

		return true
	}
	return false
}
