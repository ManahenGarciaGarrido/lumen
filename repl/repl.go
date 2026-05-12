package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ManahenGarciaGarrido/lumen/evaluator"
	"github.com/ManahenGarciaGarrido/lumen/lexer"
	"github.com/ManahenGarciaGarrido/lumen/object"
	"github.com/ManahenGarciaGarrido/lumen/parser"
)

const prompt = ">> "

const banner = `
  ██╗     ██╗   ██╗███╗   ███╗███████╗███╗   ██╗
  ██║     ██║   ██║████╗ ████║██╔════╝████╗  ██║
  ██║     ██║   ██║██╔████╔██║█████╗  ██╔██╗ ██║
  ██║     ██║   ██║██║╚██╔╝██║██╔══╝  ██║╚██╗██║
  ███████╗╚██████╔╝██║ ╚═╝ ██║███████╗██║ ╚████║
  ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝

  Lumen v0.1.0  ·  A language built from scratch in Go
  Lexer · Pratt Parser · AST · Tree-walk Evaluator

  Type .help for commands, .exit to quit.
`

// Start runs the interactive REPL.
func Start(in io.Reader, out io.Writer) {
	fmt.Fprint(out, banner)

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			fmt.Fprintln(out, "\nGoodbye!")
			return
		}

		line := scanner.Text()
		switch strings.TrimSpace(line) {
		case ".exit", ".quit":
			fmt.Fprintln(out, "Goodbye!")
			return
		case ".help":
			printHelp(out)
			continue
		case "":
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)
		prog := p.ParseProgram()

		if errs := p.Errors(); len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintln(out, e)
			}
			continue
		}

		result := evaluator.Eval(prog, env)
		if result == nil || result.Type() == "NULL" {
			continue
		}
		fmt.Fprintln(out, result.Inspect())
	}
}

func printHelp(out io.Writer) {
	fmt.Fprint(out, `
Commands:
  .exit / .quit        Exit the REPL
  .help                Show this message

Language quick-reference:
  let x = 42                    Variable declaration
  let add = fn(a, b) { a + b }  Function literal
  if (x > 0) { x } else { 0 }  Conditional expression
  [1, 2, 3]                     Array literal
  arr[0]                        Index access
  x = x + 1                     Reassignment

Built-ins:
  print(v)      len(v)     first(arr)   last(arr)
  rest(arr)     push(arr, v)   type(v)   str(v)   int(v)

`)
}
