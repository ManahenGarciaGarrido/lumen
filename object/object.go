package object

import (
	"fmt"
	"strings"

	"github.com/ManahenGarciaGarrido/lumen/ast"
)

type ObjectType string

const (
	INTEGER_OBJ  ObjectType = "INTEGER"
	STRING_OBJ   ObjectType = "STRING"
	BOOLEAN_OBJ  ObjectType = "BOOLEAN"
	NULL_OBJ     ObjectType = "NULL"
	RETURN_OBJ   ObjectType = "RETURN_VALUE"
	ERROR_OBJ    ObjectType = "ERROR"
	FUNCTION_OBJ ObjectType = "FUNCTION"
	BUILTIN_OBJ  ObjectType = "BUILTIN"
	ARRAY_OBJ    ObjectType = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer
type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// String
type String struct{ Value string }

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean
type Boolean struct{ Value bool }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue wraps a value to propagate return through block statements.
type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error
type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function holds its parameters, body, and the captured closure environment.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	return "fn(" + strings.Join(params, ", ") + ") { ... }"
}

// BuiltinFunction is the Go type for built-in functions.
type BuiltinFunction func(args ...Object) Object

type Builtin struct{ Fn BuiltinFunction }

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin" }

// Array
type Array struct{ Elements []Object }

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elems := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elems[i] = e.Inspect()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

// Environment is the variable store with optional outer scope for closures.
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Update walks the scope chain to find and mutate an existing binding.
// Returns false if the variable is not found anywhere in the chain.
func (e *Environment) Update(name string, val Object) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return false
}
