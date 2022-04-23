package main

import "github.com/koykov/entry"

type buf struct {
	idx map[string]entry.Entry64
	buf []byte
}

func (b *buf) add(s string) entry.Entry64 {
	if e, ok := b.idx[s]; ok {
		return e
	}
	var e entry.Entry64
	lo := uint32(len(b.buf))
	b.buf = append(b.buf, s...)
	hi := uint32(len(b.buf))
	e.Encode(lo, hi)
	b.idx[s] = e
	return e
}
