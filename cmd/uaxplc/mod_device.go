package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
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

type deviceTuple map[string]deviceBrand

type deviceBrand struct {
	Name   string        `yaml:"-"`
	Regex  string        `yaml:"regex"`
	Device string        `yaml:"device"`
	Model  string        `yaml:"model,omitempty"`
	Models []deviceModel `yaml:"models,omitempty"`
}

type deviceModel struct {
	Regex string `yaml:"regex"`
	Model string `yaml:"model"`
}

type deviceModule struct{}

func (m deviceModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m deviceModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "device_repo.go"
	}

	hd, _ := os.UserHomeDir()
	input = strings.ReplaceAll(input, "~", hd)

	files, err := filepath.Glob(input + "/*.yml")
	if err != nil {
		return err
	}

	var (
		bufST []deviceBrand
		bufDR []string
		bufDM []string
		bufRE []string
		buf   buf
	)

	_, _ = w.WriteString("import (\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__dr_idx = [" + strconv.Itoa(len(files)) + "][]dr{\n")

	for i := 0; i < len(files); i++ {
		bufST = bufST[:0]
		bufDR = bufDR[:0]

		var body []byte
		if body, err = ioutil.ReadFile(files[i]); err != nil {
			return
		}
		if len(body) == 0 {
			err = errors.New("nothing to parse")
			return
		}

		tuples := make(deviceTuple)
		if err = yaml.Unmarshal(body, &tuples); err != nil {
			return
		}

		for name, brand := range tuples {
			brand.Name = name
			bufST = append(bufST, brand)
		}
		sort.Slice(bufST, func(i, j int) bool {
			return bufST[i].Name < bufST[j].Name
		})

		for j := 0; j < len(bufST); j++ {
			brand := &bufST[j]
			var (
				ne entry.Entry64 // brand name index
				re = int32(-1)   // regex index
				si entry.Entry64 // string index
				sm = int32(-1)   // single model index
				me entry.Entry64 // models index
			)
			ne = buf.add(brand.Name)

			rs := brand.Regex
			if !isRegex(rs) {
				si = buf.add(brand.Regex)
			} else {
				rs = normalizeRegex(rs)
				if _, err = regexp.Compile(rs); err == nil {
					bufRE = append(bufRE, rs)
					re = int32(len(bufRE) - 1)
				} else {
					log.Printf("regexp error '%s' on '%s'", err, rs)
				}
			}

			var (
				re1 = int32(-1)   // regex index
				si1 entry.Entry64 // string index
				ne1 entry.Entry64 // model name index
			)
			if len(brand.Models) > 0 {
				meLO := uint32(len(bufDM))
				for k := 0; k < len(brand.Models); k++ {
					model := &brand.Models[k]
					rs1 := model.Regex
					if !isRegex(rs1) {
						si1 = buf.add(model.Regex)
					} else {
						rs1 = normalizeRegex(rs1)
						if _, err = regexp.Compile(rs1); err == nil {
							bufRE = append(bufRE, rs1)
							re1 = int32(len(bufRE) - 1)
						} else {
							log.Printf("regexp error '%s' on '%s'", err, rs)
						}
					}
					ne1 = buf.add(model.Model)
					bufDM = append(bufDM, fmt.Sprintf("dm{re:%s,si:%s,ne:%s}",
						hex(re1), hex(si1), hex(ne1)))
				}
				meHI := uint32(len(bufDM))
				me.Encode(meLO, meHI)
			} else if len(brand.Model) > 0 {
				ne1 = buf.add(brand.Model)
				bufDM = append(bufDM, fmt.Sprintf("dm{re:%s,si:%s,ne:%s}",
					hex(re1), hex(si1), hex(ne1)))
				sm = int32(len(bufDM)) - 1
			}

			bufDR = append(bufDR, fmt.Sprintf("dr{ne:%s,re:%s,si:%s,sm:%s,me:%s},",
				hex(ne), hex(re), hex(si), hex(sm), hex(me)))
		}

		_, _ = w.WriteString("// " + filepath.Base(files[i]) + "\n")
		_, _ = w.WriteString("{\n")
		for j := 0; j < len(bufDR); j++ {
			_, _ = w.WriteString(bufDR[j])
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString("},\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_dm = []dm{\n")
	for i := 0; i < len(bufDM); i++ {
		_, _ = w.WriteString(bufDM[i])
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_re = []*regexp.Regexp{\n")
	for i := 0; i < len(bufRE); i++ {
		_, _ = w.WriteString("regexp.MustCompile(`(?i)" + bufRE[i] + "`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_buf = []byte{\n")
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
