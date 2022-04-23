package main

import "fmt"

type vendorModule struct{}

func (m vendorModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m vendorModule) Compile(w moduleWriter, input, target string) (err error) {
	return
}
