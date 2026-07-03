package uaxpl

import (
	"math"

	"github.com/koykov/byteconv"
)

// Tokenizer walks over tokens in UA string.
type Tokenizer struct {
	pos int
}

// Next walks to next bytes token and return it with EOF flag.
func (t *Tokenizer) Next(ua []byte) ([]byte, bool) {
	if t.pos >= len(ua) {
		return nil, true // EOF
	}

	for t.pos < len(ua) && isDelimiter(ua[t.pos]) {
		t.pos++
	}
	if t.pos >= len(ua) {
		return nil, true // EOF
	}

	off := t.pos
	for t.pos < len(ua) && !isDelimiter(ua[t.pos]) {
		t.pos++
	}

	token := ua[off:t.pos]
	return token, false
}

// NextString walks to next string token and return it with EOF flag.
func (t *Tokenizer) NextString(ua string) (out string, eof bool) {
	r, eof := t.Next(byteconv.S2B(ua))
	if r != nil {
		out = byteconv.B2S(r)
	}
	return
}

// Reset resets tokenizer position.
func (t *Tokenizer) Reset() {
	t.pos = 0
}

func isDelimiter(c byte) bool {
	return tableDelim[c]
}

var tableDelim [math.MaxUint8]bool

func init() {
	tableDelim[' '] = true
	tableDelim['\t'] = true
	tableDelim['\n'] = true
	tableDelim['\r'] = true
	tableDelim['/'] = true
	tableDelim['('] = true
	tableDelim[')'] = true
	tableDelim[';'] = true
	tableDelim[','] = true
	tableDelim['.'] = true
	tableDelim['['] = true
	tableDelim[']'] = true
	tableDelim['{'] = true
	tableDelim['}'] = true
	tableDelim[':'] = true
	tableDelim['='] = true
	tableDelim['+'] = true
	tableDelim['*'] = true
}
