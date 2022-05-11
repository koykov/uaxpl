package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/koykov/entry"
	"gopkg.in/yaml.v2"
)

type osTuple struct {
	Name     string      `yaml:"name"`
	Regex    string      `yaml:"regex"`
	Version  string      `yaml:"version,omitempty"`
	Versions []osVersion `yaml:"versions,omitempty"`
}

type osVersion struct {
	Regex   string `yaml:"regex"`
	Version string `yaml:"version"`
}

type osModule struct{}

func (m osModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m osModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "os_repo.go"
	}

	hd, _ := os.UserHomeDir()
	input = strings.ReplaceAll(input, "~", hd)

	var body []byte
	if body, err = ioutil.ReadFile(input); err != nil {
		return
	}
	if len(body) == 0 {
		err = errors.New("nothing to parse")
		return
	}

	var tuples []osTuple
	if err = yaml.Unmarshal(body, &tuples); err != nil {
		return
	}

	var (
		buf   buf
		bufOR []string
		bufOV []string
		bufRE []string
	)

	for i := 0; i < len(tuples); i++ {
		tuple := &tuples[i]

		var (
			ne entry.Entry64 // name
			ni = int8(-1)
			re = int32(-1)   // regex index
			si entry.Entry64 // substring
			vi = int8(-1)    // version match index
			vs entry.Entry64 // static version
			vr entry.Entry64 // version ranges
		)

		ni, ne = m.parseInx(tuple.Name, &buf)

		rs := tuple.Regex
		if !isRegex(rs) {
			si = buf.add(tuple.Regex)
		} else {
			rs = normalizeRegex(rs)
			if _, err = regexp.Compile(rs); err == nil {
				bufRE = append(bufRE, rs)
				re = int32(len(bufRE) - 1)
			} else {
				log.Printf("regexp error '%s' on '%s'", err, rs)
			}
		}

		if len(tuple.Versions) > 0 {
			loOV := uint32(len(bufOV))
			for j := 0; j < len(tuple.Versions); j++ {
				tv := &tuple.Versions[j]
				var (
					re1 = int32(-1)   // regex index
					si1 entry.Entry64 // substring
					vi1 = int8(-1)    // version match index
					vs1 entry.Entry64 // static version
				)
				if !isRegex(tv.Regex) {
					si1 = buf.add(tv.Regex)
				} else {
					rs1 := normalizeRegex(tv.Regex)
					if _, err = regexp.Compile(rs1); err == nil {
						bufRE = append(bufRE, rs1)
						re1 = int32(len(bufRE) - 1)
					} else {
						log.Printf("regexp error '%s' on '%s'", err, rs)
					}
				}
				vi1, vs1 = m.parseInx(tv.Version, &buf)
				bufOV = append(bufOV, fmt.Sprintf("ov{re:%s,si:%s,vi:%s,vs:%s},",
					hex(re1), hex(si1), hex(vi1), hex(vs1)))
			}
			hiOV := uint32(len(bufOV))
			vr.Encode(loOV, hiOV)
		} else if len(tuple.Version) > 0 {
			vi, vs = m.parseInx(tuple.Version, &buf)
		}

		bufOR = append(bufOR, fmt.Sprintf("or{ne:%s,ni:%s,re:%s,si:%s,vi:%s,vs:%s,vr:%s},",
			hex(ne), hex(ni), hex(re), hex(si), hex(vi), hex(vs), hex(vr)))
	}

	_, _ = w.WriteString("import (\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__or_os = []or{\n")
	for i := 0; i < len(bufOR); i++ {
		_, _ = w.WriteString(bufOR[i])
		_ = w.WriteByte('\n')
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__or_ov = []ov{\n")
	for i := 0; i < len(bufOV); i++ {
		_, _ = w.WriteString(bufOV[i])
		_ = w.WriteByte('\n')
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__or_re = []*regexp.Regexp{\n")
	for i := 0; i < len(bufRE); i++ {
		_, _ = w.WriteString("regexp.MustCompile(`(?i)" + bufRE[i] + "`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__or_buf = []byte{\n")
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

func (m osModule) parseInx(ver string, buf *buf) (vi int8, vs entry.Entry64) {
	vi = -1
	if len(ver) == 0 {
		return
	}
	if ver[0] == '$' {
		vi = 1
		if i, err := strconv.ParseInt(ver[1:], 10, 64); err == nil {
			vi = int8(i)
		}
	} else {
		vs = buf.add(ver)
	}
	return
}
