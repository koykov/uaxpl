package main

import (
	"fmt"

	"github.com/koykov/entry"
)

type module interface {
	Validate(input, target string) error
	Compile(w moduleWriter, input, target string) error
}

func hex(x interface{}) string {
	switch x.(type) {
	case int8:
		i := x.(int8)
		if i < 0 {
			return fmt.Sprintf("-0x%01x", -i)
		} else {
			return fmt.Sprintf("0x%01x", i)
		}
	case int32:
		i := x.(int32)
		if i < 0 {
			return fmt.Sprintf("-0x%04x", -i)
		} else {
			return fmt.Sprintf("0x%04x", i)
		}
	case entry.Entry64:
		return fmt.Sprintf("0x%08x", x.(entry.Entry64))
	}
	return "0x0"
}
