package main

type module interface {
	Validate(input, target string) error
	Compile(w moduleWriter, input, target string) error
}
