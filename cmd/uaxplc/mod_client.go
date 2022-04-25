package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/koykov/entry"
	"gopkg.in/yaml.v2"
)

type clientTuple struct {
	Regex   string        `yaml:"regex"`
	Name    string        `yaml:"name"`
	Version string        `yaml:"version,omitempty"`
	Engine  *ClientEngine `yaml:"engine,omitempty"`
	URL     string        `yaml:"url,omitempty"`
	Type    string        `yaml:"type,omitempty"`
}

type ClientEngine struct {
	Default  string            `yaml:"default"`
	Versions map[string]string `yaml:"versions,omitempty"`
}

type clientModule struct{}

func (m clientModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m clientModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "client_repo.go"
	}

	hd, _ := os.UserHomeDir()
	input = strings.ReplaceAll(input, "~", hd)

	files, err := filepath.Glob(input + "/*.yml")
	if err != nil {
		return err
	}

	var (
		bufCR []string
		bufRE []string
		bufEF []string
		buf   buf
	)

	_, _ = w.WriteString("import (\n\"github.com/koykov/entry\"\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__cr_idx = [" + strconv.Itoa(len(files)) + "][]cr{\n")

	for i := 0; i < len(files); i++ {
		bufCR = bufCR[:0]

		var body []byte
		if body, err = ioutil.ReadFile(files[i]); err != nil {
			return
		}
		if len(body) == 0 {
			err = errors.New("nothing to parse")
			return
		}

		var tuples []clientTuple
		if err = yaml.Unmarshal(body, &tuples); err != nil {
			return
		}

		for j := 0; j < len(tuples); j++ {
			tuple := tuples[j]
			var (
				re = int32(-1)   // regex index
				si entry.Entry64 // string index
				vi = int8(-1)    // version index
				ed entry.Entry64 // default engine
				ef = int32(-1)   // engine func index
				ul entry.Entry64 // url
				tp entry.Entry64 // type string
			)

			rs := tuple.Regex
			if !isRegex(rs) {
				si = buf.add(tuple.Regex)
			} else {
				rs = normalizeRegex(rs)
				if _, err = regexp.Compile(rs); err == nil {
					bufRE = append(bufRE, rs)
					re = int32(len(bufRE) - 1)
				}
			}
			if len(tuple.Version) > 0 && tuple.Version[0] == '$' {
				n, _ := strconv.Atoi(tuple.Version[1:])
				vi = int8(n)
			}

			if tuple.Engine != nil {
				if len(tuple.Engine.Default) > 0 {
					ed = buf.add(tuple.Engine.Default)
				}
				if len(tuple.Engine.Versions) > 0 {
					// fn := "func(s string, def entry.Entry64) entry.Entry64 { return def }"
					fn := m.ef(tuple.Engine, ed, &buf)
					bufEF = append(bufEF, fn)
					ef = int32(len(bufEF) - 1)
				}
			}

			if len(tuple.URL) > 0 {
				ul = buf.add(tuple.URL)
			}
			if len(tuple.Type) > 0 {
				tp = buf.add(tuple.Type)
			}

			bufCR = append(bufCR, fmt.Sprintf("cr{re:%s,si:%s,vi:%s,ed:%s,ef:%s,ul:%s,tp:%s},",
				hex(re), hex(si), hex(vi), hex(ed), hex(ef), hex(ul), hex(tp)))
		}

		_, _ = w.WriteString("// " + filepath.Base(files[i]) + "\n")
		_, _ = w.WriteString("{\n")
		for j := 0; j < len(bufCR); j++ {
			_, _ = w.WriteString(bufCR[j])
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString("},\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__cr_re = []*regexp.Regexp{\n")
	for i := 0; i < len(bufRE); i++ {
		_, _ = w.WriteString("regexp.MustCompile(`" + bufRE[i] + "`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__cr_ef = []engFn{\n")
	for i := 0; i < len(bufEF); i++ {
		_, _ = w.WriteString(bufEF[i])
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__cr_buf = []byte{\n")
	for i := 0; i < len(buf.buf); i++ {
		if i > 0 && i%16 == 0 {
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString(fmt.Sprintf("0x%02x, ", buf.buf[i]))
	}
	_, _ = w.WriteString("\n}\n")

	_, _ = w.WriteString(")\n")

	source := w.Bytes()
	var fmtSource []byte
	if fmtSource, err = format.Source(source); err != nil {
		return
	}

	err = ioutil.WriteFile(target, fmtSource, 0644)

	return
}

func (m clientModule) ef(e *ClientEngine, def entry.Entry64, buf *buf) (out string) {
	type ve struct {
		v string
		e string
	}
	var bufVE []ve
	for v, e := range e.Versions {
		bufVE = append(bufVE, ve{v: v, e: e})
	}
	sort.Slice(bufVE, func(i, j int) bool {
		return bufVE[i].v < bufVE[j].v
	})
	out += "func(s string) entry.Entry64 {\n"
	out += "switch s {\n"
	for i := 0; i < len(bufVE); i++ {
		x := &bufVE[i]
		e1 := buf.add(x.e)
		out += "case \"" + x.v + "\":\nreturn " + fmt.Sprintf("0x%08x", e1) + "\n"
	}
	out += "default:\nreturn " + fmt.Sprintf("0x%08x", def) + "\n"
	out += "}\n}"
	return
}
