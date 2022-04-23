package main

import (
	"bytes"
	"strings"

	"github.com/koykov/fastconv"
)

var (
	regexBytes = []byte("?:.+*\\|^;()[]{}$")
)

func isRegex(s string) (b bool) {
	p := fastconv.S2B(s)
	for i := 0; i < len(regexBytes); i++ {
		if bytes.IndexByte(p, regexBytes[i]) != -1 {
			b = true
			break
		}
	}
	return
}

func normalizeRegex(s string) string {
	s = strings.ReplaceAll(s, "(?<", "(?:<")
	s = strings.ReplaceAll(s, "(?!", "(?:!")
	return s
}
