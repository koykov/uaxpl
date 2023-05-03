package uaxpl

import (
	"fmt"
	"io"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

type diff struct {
	ua  string
	buf []diffStage
}

type diffStage struct {
	key  string
	l, r string
}

func newDiff(ua string) *diff {
	d := diff{ua: ua}
	return &d
}

func (d diff) len() int {
	return len(d.buf)
}

func (d *diff) add(key, l, r string) {
	d.buf = append(d.buf, diffStage{key: key, l: bytealg.CopyStr(l), r: bytealg.CopyStr(r)})
}

func (d diff) write(w io.Writer) {
	_, _ = w.Write(fastconv.S2B(d.ua))
	_, _ = w.Write([]byte("\n"))
	for i := 0; i < len(d.buf); i++ {
		st := &d.buf[i]
		_, _ = fmt.Fprintf(w, " - %s: need '%s' got '%s'\n", st.key, st.r, st.l)
	}
}
