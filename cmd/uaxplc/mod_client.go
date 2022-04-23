package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"

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
	files, err := filepath.Glob(input + "/*.yml")
	if err != nil {
		return err
	}
	for i := 0; i < len(files); i++ {
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

		var (
			bufRE []string
			bufEF []string
			buf   buf
		)

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
			} else if _, err := regexp.Compile(normalizeRegex(rs)); err != nil {
				bufRE = append(bufRE, tuple.Regex)
				re = int32(len(bufRE) - 1)
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
					fn := "func(s string, def entry.Entry64) entry.Entry64 { return def }"
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

			_, _, _, _, _, _, _ = re, si, vi, ed, ef, ul, tp
		}
		_ = bufEF
	}

	return nil
}
