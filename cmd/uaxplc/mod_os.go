package main

import "fmt"

type osModule struct{}

func (m osModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m osModule) Compile(w moduleWriter, input, target string) (err error) {
	return
}
