package main

import (
	"errors"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"

	"github.com/koykov/entry"
	"gopkg.in/yaml.v2"
)

type vendorTuple map[string][]string

type vendorTuple1 struct {
	Name  string
	Regex []string
}

type vendorModule struct{}

func (m vendorModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m vendorModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "vendor_repo.go"
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

	tuples := make(vendorTuple)
	if err = yaml.Unmarshal(body, &tuples); err != nil {
		return
	}
	bufVT := make([]vendorTuple1, 0, len(tuples))
	for name, re := range tuples {
		bufVT = append(bufVT, vendorTuple1{
			Name:  name,
			Regex: re,
		})
	}
	sort.Slice(bufVT, func(i, j int) bool {
		return bufVT[i].Name < bufVT[j].Name
	})

	var (
		bufVR []string
		bufRE []string
		buf   buf
	)

	_, _ = w.WriteString("import (\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	for i := 0; i < len(bufVT); i++ {
		brand := &bufVT[i]

		var (
			brand64 entry.Entry64
			range64 entry.Entry64
		)

		brand64 = buf.add(brand.Name)
		lo := len(bufRE)
		for j := 0; j < len(brand.Regex); j++ {
			re := brand.Regex[j]
			bufRE = append(bufRE, re+"[^a-z0-9]+")
		}
		hi := len(bufRE)
		range64.Encode(uint32(lo), uint32(hi))

		bufVR = append(bufVR, fmt.Sprintf("vendorTuple{brand64:%s,range64:%s}",
			hex(brand64), hex(range64)))
	}

	_, _ = w.WriteString("__vr_idx = []vendorTuple{\n")
	for i := 0; i < len(bufVR); i++ {
		_, _ = w.WriteString(bufVR[i])
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__vr_re = []*regexp.Regexp{\n")
	for i := 0; i < len(bufRE); i++ {
		re := bufRE[i]
		_, _ = w.WriteString("regexp.MustCompile(`" + re + "`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__vr_buf = []byte{\n")
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
