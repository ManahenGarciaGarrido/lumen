package parser

import (
	"testing"

	"github.com/ManahenGarciaGarrido/lumen/ast"
	"github.com/ManahenGarciaGarrido/lumen/lexer"
)

func parse(input string) *ast.Program {
	return New(lexer.New(input)).ParseProgram()
}

func noErrors(t *testing.T, input string) *ast.Program {
	t.Helper()
	p := New(lexer.New(input))
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors for %q:\n%v", input, p.Errors())
	}
	return prog
}

func TestLetStatements(t *testing.T) {
	cases := []struct {
		input string
		name  string
	}{
		{"let x = 5", "x"},
		{"let foo = true", "foo"},
		{`let s = "hello"`, "s"},
		{"let result = 1 + 2", "result"},
	}
	for _, c := range cases {
		prog := noErrors(t, c.input)
		if len(prog.Statements) != 1 {
			t.Fatalf("%q: want 1 stmt, got %d", c.input, len(prog.Statements))
		}
		stmt, ok := prog.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("%q: want *LetStatement, got %T", c.input, prog.Statements[0])
		}
		if stmt.Name.Value != c.name {
			t.Errorf("%q: want name %q, got %q", c.input, c.name, stmt.Name.Value)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	prog := noErrors(t, "return 42")
	if len(prog.Statements) != 1 {
		t.Fatalf("want 1 stmt")
	}
	if _, ok := prog.Statements[0].(*ast.ReturnStatement); !ok {
		t.Fatalf("want *ReturnStatement, got %T", prog.Statements[0])
	}
}

func TestIdentifierExpression(t *testing.T) {
	prog := noErrors(t, "foobar")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	id, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("want *Identifier, got %T", stmt.Expression)
	}
	if id.Value != "foobar" {
		t.Errorf("want foobar, got %s", id.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	prog := noErrors(t, "42")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	lit, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("want *IntegerLiteral, got %T", stmt.Expression)
	}
	if lit.Value != 42 {
		t.Errorf("want 42, got %d", lit.Value)
	}
}

func TestBooleanLiteral(t *testing.T) {
	for _, c := range []struct {
		input string
		want  bool
	}{{"true", true}, {"false", false}} {
		prog := noErrors(t, c.input)
		stmt := prog.Statements[0].(*ast.ExpressionStatement)
		b, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("want *BooleanLiteral, got %T", stmt.Expression)
		}
		if b.Value != c.want {
			t.Errorf("want %v, got %v", c.want, b.Value)
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	cases := []struct {
		input string
		op    string
	}{
		{"!true", "!"},
		{"-5", "-"},
		{"!false", "!"},
	}
	for _, c := range cases {
		prog := noErrors(t, c.input)
		stmt := prog.Statements[0].(*ast.ExpressionStatement)
		pe, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("%q: want *PrefixExpression, got %T", c.input, stmt.Expression)
		}
		if pe.Operator != c.op {
			t.Errorf("%q: want op %q, got %q", c.input, c.op, pe.Operator)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	cases := []struct {
		input string
		op    string
	}{
		{"1 + 2", "+"},
		{"3 - 1", "-"},
		{"2 * 4", "*"},
		{"10 / 2", "/"},
		{"5 % 3", "%"},
		{"1 == 1", "=="},
		{"1 != 2", "!="},
		{"1 < 2", "<"},
		{"2 > 1", ">"},
		{"1 <= 1", "<="},
		{"2 >= 2", ">="},
		{"true && false", "&&"},
		{"true || false", "||"},
	}
	for _, c := range cases {
		prog := noErrors(t, c.input)
		stmt := prog.Statements[0].(*ast.ExpressionStatement)
		ie, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("%q: want *InfixExpression, got %T", c.input, stmt.Expression)
		}
		if ie.Operator != c.op {
			t.Errorf("%q: want op %q, got %q", c.input, c.op, ie.Operator)
		}
	}
}

func TestGroupedExpression(t *testing.T) {
	prog := noErrors(t, "(1 + 2) * 3")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	ie, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("want *InfixExpression, got %T", stmt.Expression)
	}
	if ie.Operator != "*" {
		t.Errorf("want *, got %s", ie.Operator)
	}
}

func TestIfExpression(t *testing.T) {
	prog := noErrors(t, "if (x > 5) { x } else { 0 }")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	ie, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("want *IfExpression, got %T", stmt.Expression)
	}
	if ie.Alternative == nil {
		t.Fatal("want alternative, got nil")
	}
}

func TestIfWithoutElse(t *testing.T) {
	prog := noErrors(t, "if (x > 0) { x }")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	ie := stmt.Expression.(*ast.IfExpression)
	if ie.Alternative != nil {
		t.Error("want nil alternative")
	}
}

func TestFunctionLiteral(t *testing.T) {
	prog := noErrors(t, "fn(x, y) { x + y }")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("want *FunctionLiteral, got %T", stmt.Expression)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("want 2 params, got %d", len(fn.Parameters))
	}
	if fn.Parameters[0].Value != "x" || fn.Parameters[1].Value != "y" {
		t.Errorf("wrong params: %v", fn.Parameters)
	}
}

func TestFunctionNoParams(t *testing.T) {
	prog := noErrors(t, "fn() { 42 }")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	fn := stmt.Expression.(*ast.FunctionLiteral)
	if len(fn.Parameters) != 0 {
		t.Errorf("want 0 params, got %d", len(fn.Parameters))
	}
}

func TestCallExpression(t *testing.T) {
	prog := noErrors(t, "add(1, 2 * 3)")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("want *CallExpression, got %T", stmt.Expression)
	}
	if len(call.Arguments) != 2 {
		t.Fatalf("want 2 args, got %d", len(call.Arguments))
	}
}

func TestArrayLiteral(t *testing.T) {
	prog := noErrors(t, "[1, 2, 3]")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("want *ArrayLiteral, got %T", stmt.Expression)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("want 3 elements, got %d", len(arr.Elements))
	}
}

func TestIndexExpression(t *testing.T) {
	prog := noErrors(t, "arr[0]")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	if _, ok := stmt.Expression.(*ast.IndexExpression); !ok {
		t.Fatalf("want *IndexExpression, got %T", stmt.Expression)
	}
}

func TestAssignExpression(t *testing.T) {
	prog := noErrors(t, "x = x + 1")
	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	ae, ok := stmt.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("want *AssignExpression, got %T", stmt.Expression)
	}
	if ae.Name.Value != "x" {
		t.Errorf("want x, got %s", ae.Name.Value)
	}
}

func TestParserErrors(t *testing.T) {
	bad := []string{
		"let = 5",
		"let x 5",
		"fn(,) {}",
	}
	for _, input := range bad {
		p := New(lexer.New(input))
		p.ParseProgram()
		if len(p.Errors()) == 0 {
			t.Errorf("%q: expected parse errors, got none", input)
		}
	}
}
