package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/koykov/hash/fnv"
	"gopkg.in/yaml.v2"
)

type hintTuple map[string]string

type hintModule struct{}

func (m hintModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m hintModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "hint_repo.go"
	}

	hd, _ := os.UserHomeDir()
	input = strings.ReplaceAll(input, "~", hd)

	files, err := filepath.Glob(input + "/*.yml")
	if err != nil {
		return err
	}

	var (
		bufST []string
		uniqK = make(map[string]struct{})
		buf   buf
	)

	_, _ = w.WriteString("import (\n\"github.com/koykov/entry\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__hr_idx = map[uint64]entry.Entry64{\n")

	for i := 0; i < len(files); i++ {
		var body []byte
		if body, err = ioutil.ReadFile(files[i]); err != nil {
			return
		}
		if len(body) == 0 {
			err = errors.New("nothing to parse")
			return
		}

		var tuples hintTuple
		if err = yaml.Unmarshal(body, &tuples); err != nil {
			return
		}

		for k := range tuples {
			bufST = append(bufST, k)
		}
		sort.Strings(bufST)

		_, _ = w.WriteString("// " + filepath.Base(files[i]) + "\n")
		for j := 0; j < len(bufST); j++ {
			k := bufST[j]
			if _, ok := uniqK[k]; ok {
				continue
			}
			uniqK[k] = struct{}{}
			h := fnv.Hash64String(k)
			e := buf.add(tuples[k])
			_, _ = w.WriteString(fmt.Sprintf("%s:%s,\n", hex(h), hex(e)))
		}
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__hr_buf = []byte{\n")
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
