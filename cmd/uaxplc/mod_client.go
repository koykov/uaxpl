package main

import (
	"errors"
	"fmt"
	"go/format"
	"log"
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

var (
	// Regexp replaces
	reReplClient = map[string]string{
		`(?:Go-http-client|Go )/?(?:(\d+[\.\d]+))?(?: package http)?`: `(?:Go-http-client)/?(?:(\d+[\.\d]+))?(?: package http)?`,
	}
)

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
		bufCR  []string
		bufRE  []string
		bufEF  []string
		idxE   = make(map[string]struct{})
		bufERE []string
		buf    buf
	)

	_, _ = w.WriteString("import (\n\"github.com/koykov/entry\"\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__cr_idx = [" + strconv.Itoa(len(files)) + "][]clientTuple{\n")

	for i := 0; i < len(files); i++ {
		bufCR = bufCR[:0]

		var body []byte
		if body, err = os.ReadFile(files[i]); err != nil {
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
				matchRI   = int32(-1)   // regex index
				match64   entry.Entry64 // string index
				browser64 entry.Entry64 // browser name index
				version64 entry.Entry64 // browser version index
				browserVI = int8(-1)    // version index
				engine64  entry.Entry64 // default engine
				engineFI  = int32(-1)   // engine func index
				url64     entry.Entry64 // url
				type64    entry.Entry64 // type string
			)

			rs := tuple.Regex
			if !isRegex(rs) {
				match64 = buf.add(tuple.Regex)
			} else {
				rs = normalizeRegex(rs)
				if rs1, ok := reReplClient[rs]; ok {
					rs = rs1
				}
				if _, err = regexp.Compile(rs); err == nil {
					bufRE = append(bufRE, rs)
					matchRI = int32(len(bufRE) - 1)
				} else {
					log.Printf("regexp error '%s' on '%s'", err, rs)
				}
			}
			if len(tuple.Name) > 0 {
				browser64 = buf.add(tuple.Name)
			}
			if len(tuple.Version) > 0 {
				if tuple.Version[0] == '$' {
					n, _ := strconv.Atoi(tuple.Version[1:])
					browserVI = int8(n)
				} else {
					version64 = buf.add(tuple.Version)
				}
			}

			if tuple.Engine != nil {
				if len(tuple.Engine.Default) > 0 {
					engine64 = buf.add(tuple.Engine.Default)
					idxE[tuple.Engine.Default] = struct{}{}
				}
				if len(tuple.Engine.Versions) > 0 {
					fn := m.ef(tuple.Engine, engine64, &buf)
					bufEF = append(bufEF, fn)
					engineFI = int32(len(bufEF) - 1)
					for _, e := range tuple.Engine.Versions {
						idxE[e] = struct{}{}
					}
				}
			}

			if len(tuple.URL) > 0 {
				url64 = buf.add(tuple.URL)
			}
			if len(tuple.Type) > 0 {
				type64 = buf.add(tuple.Type)
			}

			bufCR = append(bufCR, fmt.Sprintf("clientTuple{matchRI:%s,match64:%s,browser64:%s,version64:%s,browserVI:%s,engine64:%s,engineFI:%s,url64:%s,type64:%s},",
				hex(matchRI), hex(match64), hex(browser64), hex(version64), hex(browserVI), hex(engine64), hex(engineFI), hex(url64), hex(type64)))
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
		_, _ = w.WriteString("regexp.MustCompile(`(?i)(?:^|[^A-Z0-9\\-_]|[^A-Z0-9\\-]_|sprd-|MZ-)(?:" + bufRE[i] + ")`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__cr_ef = []engineFn{\n")
	for i := 0; i < len(bufEF); i++ {
		_, _ = w.WriteString(bufEF[i])
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__cr_ev = map[entry.Entry64]int32{\n")
	if len(idxE) > 0 {
		var bs []string
		for e := range idxE {
			if len(e) == 0 {
				continue
			}
			bs = append(bs, e)
		}
		sort.Strings(bs)
		for i := 0; i < len(bs); i++ {
			var re string
			en := bs[i]
			if en == "Gecko" {
				re = `[ ](?:rv[: ]([0-9\.]+)).*gecko/[0-9.]+`
			} else {
				if en == "Blink" {
					en = "Chrome"
				}
				re = fmt.Sprintf(`%s\s*[/\s]\s*(\d+(?:.\d+)+)`, en)
			}
			e := buf.add(bs[i])
			bufERE = append(bufERE, re)
			_, _ = w.WriteString(hex(e))
			_ = w.WriteByte(':')
			_, _ = w.WriteString(hex(int32(len(bufERE) - 1)))
			_, _ = w.WriteString(",\n")
		}
	}
	_, _ = w.WriteString("}\n")
	_, _ = w.WriteString("__cr_evre = []*regexp.Regexp{\n")
	for i := 0; i < len(bufERE); i++ {
		_, _ = w.WriteString("regexp.MustCompile(`(?i)" + bufERE[i] + "`),\n")
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

	err = os.WriteFile(target, fmtSource, 0644)

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
	for i := 0; i < len(bufVE); i++ {
		x := &bufVE[i]
		e1 := buf.add(x.e)
		out += "if s>=\"" + x.v + "\"{return " + fmt.Sprintf("0x%08x", e1) + "}\n"
	}
	out += "return " + fmt.Sprintf("0x%08x", def) + "\n"
	out += "}"
	return
}
