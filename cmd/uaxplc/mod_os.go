package main

import (
	"errors"
	"fmt"
	"go/format"
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

var (
	// Regexp replaces
	reReplOS = map[string]string{
		`(?:iPhone ?OS|iOS(?: Version)?)(?:/|; |,)(\d+[.\d]+)`:                           `(?:iPhone ?OS|[\s(]iOS(?: Version)?)(?:/|; |,)(\d+[\.\d]+)`,
		`Android-(\d+[.\d]*);`:                                                           `Android[\s\-](\d+[.\d]*);`,
		`(?:Android API \d+|\d+/tclwebkit(?:\d+[.\d]*)|(?:Android/\d{2}|Android \d{2}))`: `(?:Android API \d+|\d+/tclwebkit(?:\d+[.\d]*)|(?:[^_]Android/\d{2}|Android-\d{2}))`,
	}
)

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
	if body, err = os.ReadFile(input); err != nil {
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
			name64     entry.Entry64 // name
			nameSI     = int8(-1)
			matchRI    = int32(-1)   // regex index
			match64    entry.Entry64 // substring
			versionSI  = int8(-1)    // version match index
			version64  entry.Entry64 // static version
			versions64 entry.Entry64 // version ranges
		)

		nameSI, name64 = m.parseInx(tuple.Name, &buf)

		rs := tuple.Regex
		if !isRegex(rs) {
			match64 = buf.add(tuple.Regex)
		} else {
			rs = normalizeRegex(rs)
			if rs1, ok := reReplOS[rs]; ok {
				rs = rs1
			}
			if _, err = regexp.Compile(rs); err == nil {
				bufRE = append(bufRE, rs)
				matchRI = int32(len(bufRE) - 1)
			} else {
				log.Printf("regexp error '%s' on '%s'", err, rs)
			}
		}

		if len(tuple.Versions) > 0 {
			loOV := uint32(len(bufOV))
			for j := 0; j < len(tuple.Versions); j++ {
				tv := &tuple.Versions[j]
				var (
					matchRI1   = int32(-1)   // regex index
					match641   entry.Entry64 // substring
					versionSI1 = int8(-1)    // version match index
					version641 entry.Entry64 // static version
				)
				if !isRegex(tv.Regex) {
					match641 = buf.add(tv.Regex)
				} else {
					rs1 := normalizeRegex(tv.Regex)
					if rs1, ok := reReplOS[rs]; ok {
						rs = rs1
					}
					if _, err = regexp.Compile(rs1); err == nil {
						bufRE = append(bufRE, rs1)
						matchRI1 = int32(len(bufRE) - 1)
					} else {
						log.Printf("regexp error '%s' on '%s'", err, rs)
					}
				}
				versionSI1, version641 = m.parseInx(tv.Version, &buf)
				bufOV = append(bufOV, fmt.Sprintf("osVersionTuple{matchRI:%s,match64:%s,versionSI:%s,version64:%s},",
					hex(matchRI1), hex(match641), hex(versionSI1), hex(version641)))
			}
			hiOV := uint32(len(bufOV))
			versions64.Encode(loOV, hiOV)
		} else if len(tuple.Version) > 0 {
			versionSI, version64 = m.parseInx(tuple.Version, &buf)
		}

		bufOR = append(bufOR, fmt.Sprintf("osTuple{name64:%s,nameSI:%s,matchRI:%s,match64:%s,versionSI:%s,version64:%s,versions64:%s},",
			hex(name64), hex(nameSI), hex(matchRI), hex(match64), hex(versionSI), hex(version64), hex(versions64)))
	}

	_, _ = w.WriteString("import (\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__or_os = []osTuple{\n")
	for i := 0; i < len(bufOR); i++ {
		_, _ = w.WriteString(bufOR[i])
		_ = w.WriteByte('\n')
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__or_ov = []osVersionTuple{\n")
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

	err = os.WriteFile(target, fmtSource, 0644)

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
