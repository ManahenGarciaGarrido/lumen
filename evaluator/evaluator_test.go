package evaluator

import (
	"testing"

	"github.com/ManahenGarciaGarrido/lumen/lexer"
	"github.com/ManahenGarciaGarrido/lumen/object"
	"github.com/ManahenGarciaGarrido/lumen/parser"
)

func eval(input string) object.Object {
	prog := parser.New(lexer.New(input)).ParseProgram()
	return Eval(prog, object.NewEnvironment())
}

func asInt(t *testing.T, obj object.Object, want int64) {
	t.Helper()
	i, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("want *Integer(%d), got %T: %s", want, obj, obj.Inspect())
	}
	if i.Value != want {
		t.Errorf("want %d, got %d", want, i.Value)
	}
}

func asBool(t *testing.T, obj object.Object, want bool) {
	t.Helper()
	b, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("want *Boolean(%v), got %T", want, obj)
	}
	if b.Value != want {
		t.Errorf("want %v, got %v", want, b.Value)
	}
}

func asStr(t *testing.T, obj object.Object, want string) {
	t.Helper()
	s, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("want *String(%q), got %T", want, obj)
	}
	if s.Value != want {
		t.Errorf("want %q, got %q", want, s.Value)
	}
}

func asError(t *testing.T, obj object.Object) *object.Error {
	t.Helper()
	e, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("want *Error, got %T: %s", obj, obj.Inspect())
	}
	return e
}

// ── Integer arithmetic ────────────────────────────────────────────────────────

func TestIntegerArithmetic(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"0", 0},
		{"5 + 5 + 5", 15},
		{"10 - 3", 7},
		{"4 * 3", 12},
		{"10 / 2", 5},
		{"10 % 3", 1},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"100 / 10 / 2", 5},
		{"-5", -5},
		{"-10 + 5", -5},
		{"--5", 5},
	}
	for _, c := range cases {
		asInt(t, eval(c.input), c.want)
	}
}

// ── Boolean expressions ───────────────────────────────────────────────────────

func TestBooleanExpressions(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 <= 1", true},
		{"2 >= 2", true},
		{"1 == 1", true},
		{"1 != 2", true},
		{"1 == 2", false},
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"1 < 2 && 3 < 4", true},
		{"1 > 2 && 3 < 4", false},
		{"1 > 2 || 3 < 4", true},
		{"1 > 2 || 3 > 4", false},
	}
	for _, c := range cases {
		asBool(t, eval(c.input), c.want)
	}
}

// ── Variables ─────────────────────────────────────────────────────────────────

func TestLetStatements(t *testing.T) {
	asInt(t, eval("let x = 5; x"), 5)
	asInt(t, eval("let x = 5; let y = 3; x + y"), 8)
	asInt(t, eval("let x = 5 * 5; x"), 25)
}

func TestAssignment(t *testing.T) {
	asInt(t, eval("let x = 1\nx = x + 1\nx = x + 1\nx"), 3)
	asInt(t, eval("let x = 10\nx = x * 2\nx"), 20)
}

// ── Conditionals ──────────────────────────────────────────────────────────────

func TestIfElse(t *testing.T) {
	asInt(t, eval("if (true) { 10 }"), 10)
	asInt(t, eval("if (false) { 10 } else { 20 }"), 20)
	asInt(t, eval("if (1 < 2) { 3 } else { 4 }"), 3)
	asInt(t, eval("if (1 > 2) { 3 } else { 4 }"), 4)
	asInt(t, eval("if (1 == 1) { 99 }"), 99)
}

func TestIfNoElseReturnsNull(t *testing.T) {
	result := eval("if (false) { 10 }")
	if _, ok := result.(*object.Null); !ok {
		t.Errorf("want Null, got %T", result)
	}
}

// ── Return ────────────────────────────────────────────────────────────────────

func TestReturnStatements(t *testing.T) {
	asInt(t, eval("return 5; 10"), 5)
	asInt(t, eval("return 2 * 5; return 1"), 10)
	asInt(t, eval("if (true) { return 7 }\nreturn 0"), 7)
}

// ── Functions ─────────────────────────────────────────────────────────────────

func TestFunctions(t *testing.T) {
	asInt(t, eval("let double = fn(x) { x * 2 }; double(5)"), 10)
	asInt(t, eval("let add = fn(a, b) { a + b }; add(3, 4)"), 7)
	asInt(t, eval("let sq = fn(x) { x * x }; sq(9)"), 81)
	asInt(t, eval("fn(x) { x + 1 }(10)"), 11)
}

func TestHigherOrderFunctions(t *testing.T) {
	input := `
let apply = fn(f, x) { f(x) }
let double = fn(x) { x * 2 }
apply(double, 5)`
	asInt(t, eval(input), 10)
}

// ── Closures ──────────────────────────────────────────────────────────────────

func TestClosures(t *testing.T) {
	input := `
let makeAdder = fn(x) {
  fn(y) { x + y }
}
let add5 = makeAdder(5)
add5(3)`
	asInt(t, eval(input), 8)
}

func TestClosureMutation(t *testing.T) {
	input := `
let makeCounter = fn() {
  let count = 0
  return fn() {
    count = count + 1
    return count
  }
}
let next = makeCounter()
next()
next()
next()`
	asInt(t, eval(input), 3)
}

// ── Recursion ─────────────────────────────────────────────────────────────────

func TestRecursionFib(t *testing.T) {
	input := `
let fib = fn(n) {
  if (n < 2) { return n }
  return fib(n - 1) + fib(n - 2)
}
fib(10)`
	asInt(t, eval(input), 55)
}

func TestRecursionFactorial(t *testing.T) {
	input := `
let fact = fn(n) {
  if (n <= 1) { return 1 }
  return n * fact(n - 1)
}
fact(6)`
	asInt(t, eval(input), 720)
}

// ── Strings ───────────────────────────────────────────────────────────────────

func TestStringConcatenation(t *testing.T) {
	asStr(t, eval(`"hello" + ", " + "world"`), "hello, world")
	asStr(t, eval(`"a" + "b" + "c"`), "abc")
}

func TestStringEquality(t *testing.T) {
	asBool(t, eval(`"abc" == "abc"`), true)
	asBool(t, eval(`"abc" != "xyz"`), true)
	asBool(t, eval(`"abc" == "xyz"`), false)
}

// ── Arrays ────────────────────────────────────────────────────────────────────

func TestArrayIndex(t *testing.T) {
	asInt(t, eval("[1, 2, 3][0]"), 1)
	asInt(t, eval("[1, 2, 3][1]"), 2)
	asInt(t, eval("[1, 2, 3][2]"), 3)
	asInt(t, eval("let a = [10, 20, 30]; a[2]"), 30)
}

func TestArrayOutOfBoundsIsNull(t *testing.T) {
	result := eval("[1, 2, 3][10]")
	if _, ok := result.(*object.Null); !ok {
		t.Errorf("want Null, got %T", result)
	}
}

// ── Built-ins ─────────────────────────────────────────────────────────────────

func TestBuiltinLen(t *testing.T) {
	asInt(t, eval(`len("hello")`), 5)
	asInt(t, eval(`len("")`), 0)
	asInt(t, eval("len([1, 2, 3])"), 3)
	asInt(t, eval("len([])"), 0)
}

func TestBuiltinFirst(t *testing.T) {
	asInt(t, eval("first([10, 20, 30])"), 10)
	result := eval("first([])")
	if _, ok := result.(*object.Null); !ok {
		t.Errorf("want Null for first([]), got %T", result)
	}
}

func TestBuiltinLast(t *testing.T) {
	asInt(t, eval("last([10, 20, 30])"), 30)
}

func TestBuiltinRest(t *testing.T) {
	result := eval("rest([1, 2, 3])")
	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("want Array, got %T", result)
	}
	if len(arr.Elements) != 2 {
		t.Errorf("want 2 elements, got %d", len(arr.Elements))
	}
}

func TestBuiltinPush(t *testing.T) {
	result := eval("push([1, 2], 3)")
	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("want Array, got %T", result)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("want 3 elements, got %d", len(arr.Elements))
	}
	asInt(t, arr.Elements[2], 3)
}

func TestBuiltinType(t *testing.T) {
	asStr(t, eval("type(42)"), "integer")
	asStr(t, eval(`type("hi")`), "string")
	asStr(t, eval("type(true)"), "boolean")
	asStr(t, eval("type([])"), "array")
}

func TestBuiltinStr(t *testing.T) {
	asStr(t, eval("str(42)"), "42")
	asStr(t, eval("str(true)"), "true")
}

func TestBuiltinInt(t *testing.T) {
	asInt(t, eval(`int("42")`), 42)
	asInt(t, eval("int(true)"), 1)
	asInt(t, eval("int(false)"), 0)
}

// ── Error handling ────────────────────────────────────────────────────────────

func TestErrors(t *testing.T) {
	cases := []string{
		"foobar",
		"5 + true",
		"true + false",
		"1 / 0",
		"1 % 0",
		`len(1, 2)`,
	}
	for _, c := range cases {
		result := eval(c)
		if _, ok := result.(*object.Error); !ok {
			t.Errorf("%q: want Error, got %T: %s", c, result, result.Inspect())
		}
	}
}

func TestErrorsDoNotPropagate(t *testing.T) {
	// An error in arg evaluation should short-circuit
	result := eval("let x = 1/0; x")
	asError(t, result)
}

func TestDivisionByZero(t *testing.T) {
	e := asError(t, eval("5 / 0"))
	if e.Message == "" {
		t.Error("want non-empty error message")
	}
}
