package main

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/koykov/fastconv"
)

var (
	regexBytes            = []byte("?:.+*\\|^;()[]{}$")
	reNegativeLookbehind  = regexp.MustCompile(`\(\?<!([^)]+)\)`)
	reNegativeLookaheadWA = regexp.MustCompile(`\(\?!\.\*([^)]+)\)`)
	reNegativeLookahead   = regexp.MustCompile(`\(\?!([^)]+)\)`)
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
	if reNegativeLookbehind.MatchString(s) {
		for {
			bo, bc := 1, 0
			io := strings.Index(s, "(?<") + 1
			if io == 0 {
				break
			}
			ic := io
			for i := io; i < len(s); i++ {
				if s[i] == '(' {
					bo++
				}
				if s[i] == ')' && s[i-1] != '\\' {
					bc++
				}
				if bo == bc {
					ic = i
					break
				}
			}
			ss := s[io-1 : ic+1]
			// sr := "[^" + ss[4:len(ss)-1] + "]*"
			// sr = strings.ReplaceAll(sr, "-", "\\-")
			sr := ""
			s = strings.Replace(s, ss, sr, 1)
		}
	}
	if reNegativeLookaheadWA.MatchString(s) {
		bo, bc := 1, 0
		io := strings.Index(s, "(?!.*") + 1
		ic := io
		for i := io; i < len(s); i++ {
			if s[i] == '(' {
				bo++
			}
			if s[i] == ')' && s[i-1] != '\\' {
				bc++
			}
			if bo == bc {
				ic = i
				break
			}
		}
		ss := s[io-1 : ic+1]
		// sr := ".*[^" + ss[5:len(ss)-1] + "]*"
		// sr = strings.ReplaceAll(sr, "-", "\\-")
		sr := ""
		s = strings.Replace(s, ss, sr, 1)
	}
	if reNegativeLookahead.MatchString(s) {
		for {
			bo, bc := 1, 0
			io := strings.Index(s, "(?!") + 1
			if io == 0 {
				break
			}
			ic := io
			for i := io; i < len(s); i++ {
				if s[i] == '(' {
					bo++
				}
				if s[i] == ')' && s[i-1] != '\\' {
					bc++
				}
				if bo == bc {
					ic = i
					break
				}
			}
			ss := s[io-1 : ic+1]
			// sr := "[^" + ss[3:len(ss)-1] + "]*"
			// sr = strings.ReplaceAll(sr, "-", "\\-")
			sr := ""
			s = strings.Replace(s, ss, sr, 1)
		}
	}
	return s
}
