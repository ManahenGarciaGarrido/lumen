package main

import (
	"fmt"
	"os"

	"github.com/ManahenGarciaGarrido/lumen/evaluator"
	"github.com/ManahenGarciaGarrido/lumen/lexer"
	"github.com/ManahenGarciaGarrido/lumen/object"
	"github.com/ManahenGarciaGarrido/lumen/parser"
	"github.com/ManahenGarciaGarrido/lumen/repl"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		repl.Start(os.Stdin, os.Stdout)
	case 1:
		runFile(args[0])
	default:
		fmt.Fprintln(os.Stderr, "Usage: lumen [file.lum]")
		os.Exit(1)
	}
}

func runFile(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read %q: %v\n", path, err)
		os.Exit(1)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	prog := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(prog, env)
	if result != nil && result.Type() == "ERROR" {
		fmt.Fprintln(os.Stderr, result.Inspect())
		os.Exit(1)
	}
}
