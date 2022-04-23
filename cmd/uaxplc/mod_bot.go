package main

import "fmt"

type botModule struct{}

func (m botModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m botModule) Compile(w moduleWriter, input, target string) (err error) {
	return
}
