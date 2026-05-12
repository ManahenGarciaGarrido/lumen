# Lumen — A Programming Language Built from Scratch

> Designed and implemented from zero: Lexer → Pratt Parser → AST → Tree-walk Evaluator

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)
![Tests](https://img.shields.io/badge/tests-60%2B%20passing-98C379?style=flat)
![License](https://img.shields.io/badge/License-MIT-7C6AF7?style=flat)
[![CI](https://github.com/ManahenGarciaGarrido/lumen/actions/workflows/ci.yml/badge.svg)](https://github.com/ManahenGarciaGarrido/lumen/actions)

---

## What is this?

Lumen is a dynamically-typed scripting language implemented entirely from scratch in Go — no parser generators, no external libraries. It walks the canonical compiler-construction pipeline:

```
source code  →  Lexer  →  Token stream  →  Parser  →  AST  →  Evaluator  →  values
```

The goal is to demonstrate deep understanding of how programming languages actually work, not just how to use them.

---

## Language features

| Feature | Example |
|---|---|
| Variables | `let x = 42` |
| Functions (first-class) | `let add = fn(a, b) { a + b }` |
| Closures | `let counter = makeCounter()` |
| Recursion | `fib(n - 1) + fib(n - 2)` |
| Conditionals | `if (x > 0) { x } else { 0 }` |
| Reassignment | `x = x + 1` |
| Arrays | `[1, 2, 3]` |
| String concat | `"Hello, " + name + "!"` |
| Operators | `+ - * / % == != < > <= >= && \|\|` |
| Built-ins | `print len first last rest push type str int` |

---

## Architecture

```
lumen/
├── token/        Token types and the Token struct (Type, Literal, Line, Col)
├── lexer/        Single-pass character scanner → token stream
├── ast/          All AST node types (Program, LetStatement, IfExpression …)
├── parser/       Recursive descent + Pratt parsing for operator precedence
├── object/       Runtime value types + Environment (lexical scope chain)
├── evaluator/    Tree-walk interpreter + built-in functions
├── repl/         Interactive REPL with banner and command history
├── playground/   Self-contained browser playground (HTML/CSS/JS)
└── examples/     Sample .lum programs
```

### Key design decisions

**Pratt Parser** — each token type carries a prefix and/or infix parse function. Precedence is handled by a simple integer table, enabling correct `2 + 3 * 4 = 14` without grammar ambiguity.

**Environment chain** — functions close over their definition-time environment. Calling a function creates a new enclosed environment pointing to the closure, not the call site. This gives correct lexical scoping for free.

**Tree-walk evaluation** — no bytecode, no compilation step. The evaluator pattern-matches on AST node types via Go's type switch and recurses. Simple to understand, easy to extend.

---

## Getting started

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/ManahenGarciaGarrido/lumen.git
cd lumen
go build -o lumen .
```

### REPL

```bash
./lumen
```

```
  ██╗     ██╗   ██╗███╗   ███╗███████╗███╗   ██╗
  ...

  Lumen v0.1.0  ·  A language built from scratch in Go

>> let fib = fn(n) { if (n < 2) { return n }  return fib(n-1) + fib(n-2) }
>> fib(10)
55
>> .exit
```

### Run a file

```bash
./lumen examples/fibonacci.lum
./lumen examples/closures.lum
./lumen examples/arrays.lum
```

### Run tests

```bash
go test ./... -v
```

---

## Examples

### Fibonacci (recursion)
```lumen
let fib = fn(n) {
  if (n < 2) { return n }
  return fib(n - 1) + fib(n - 2)
}

print(fib(10))  // 55
```

### Counter (closures + mutation)
```lumen
let makeCounter = fn() {
  let count = 0
  return fn() {
    count = count + 1
    return count
  }
}

let next = makeCounter()
print(next())  // 1
print(next())  // 2
print(next())  // 3
```

### Higher-order functions
```lumen
let apply = fn(f, x) { f(x) }
let double = fn(x) { x * 2 }

print(apply(double, 21))  // 42
```

---

## Browser playground

Open `playground/index.html` in any browser — no build step required. Includes live syntax highlighting, token table, and interactive AST visualizer.

---

## Roadmap

- [ ] Bytecode compiler + stack-based VM
- [ ] Hash maps `{"key": value}`
- [ ] `for` / `while` loops
- [ ] Standard library (`io`, `math`)
- [ ] Module system

---

## References

Thorsten Ball — *Writing An Interpreter In Go* (the canonical reference for this type of project).

---

## License

MIT © Manahen García Garrido
