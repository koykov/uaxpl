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

func (c *Ctx) parseClient() bool {
	if c.maskClientType&ClientTypeFeedReader != 0 {
		if c.evalClient(cpFeedReader) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypeFeedReader
			return true
		}
	}
	if c.maskClientType&ClientTypeMobileApp != 0 {
		if c.evalClient(cpMobileApp) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypeMobileApp
			return true
		}
	}
	if c.maskClientType&ClientTypeMediaPlayer != 0 {
		if c.evalClient(cpMediaPlayer) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypeMediaPlayer
			return true
		}
	}
	if c.maskClientType&ClientTypePIM != 0 {
		if c.evalClient(cpPIM) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypePIM
			return true
		}
	}
	if c.maskClientType&ClientTypeLibrary != 0 {
		if c.evalClient(cpLibrary) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypeLibrary
			return true
		}
	}
	if c.maskClientType&ClientTypeBrowser != 0 {
		if c.evalClient(cpBrowser) {
			c.SetBit(flagClientDetect, true)
			c.clientType = ClientTypeBrowser
			return true
		}
	}
	return false
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
