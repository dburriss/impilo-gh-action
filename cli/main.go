package main

import (
	"dburriss/impilo_gh/input"
	"fmt"
	"os"
)

//go:generate go run ./input/gen.go

func main() {
	input := input.NewActionInput(os.Args[1:])
	fmt.Print(input)
}
