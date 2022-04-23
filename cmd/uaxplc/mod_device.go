package main

import "fmt"

type deviceModule struct{}

func (m deviceModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m deviceModule) Compile(w moduleWriter, input, target string) (err error) {
	return
}
