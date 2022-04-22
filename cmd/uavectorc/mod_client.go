package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

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

		for j := 0; j < len(tuples); j++ {
			fmt.Println(tuples[j].Regex)
		}
	}

	return nil
}
