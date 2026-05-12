package lexer

import (
	"testing"

	"github.com/ManahenGarciaGarrido/lumen/token"
)

func TestNextToken(t *testing.T) {
	input := `let x = 5
let add = fn(a, b) {
  return a + b
}
print(add(x, 10))
let s = "hello world"
if (x > 3) { true } else { false }`

	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.LET, "let"}, {token.IDENT, "x"}, {token.ASSIGN, "="}, {token.INT, "5"},
		{token.LET, "let"}, {token.IDENT, "add"}, {token.ASSIGN, "="}, {token.FUNCTION, "fn"},
		{token.LPAREN, "("}, {token.IDENT, "a"}, {token.COMMA, ","}, {token.IDENT, "b"}, {token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"}, {token.IDENT, "a"}, {token.PLUS, "+"}, {token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.IDENT, "print"}, {token.LPAREN, "("},
		{token.IDENT, "add"}, {token.LPAREN, "("}, {token.IDENT, "x"}, {token.COMMA, ","}, {token.INT, "10"}, {token.RPAREN, ")"}, {token.RPAREN, ")"},
		{token.LET, "let"}, {token.IDENT, "s"}, {token.ASSIGN, "="}, {token.STRING, "hello world"},
		{token.IF, "if"}, {token.LPAREN, "("}, {token.IDENT, "x"}, {token.GT, ">"}, {token.INT, "3"}, {token.RPAREN, ")"},
		{token.LBRACE, "{"}, {token.TRUE, "true"}, {token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"}, {token.FALSE, "false"}, {token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType {
			t.Errorf("[%d] type: want %q got %q (literal=%q)", i, tt.wantType, tok.Type, tok.Literal)
		}
		if tok.Literal != tt.wantLiteral {
			t.Errorf("[%d] literal: want %q got %q", i, tt.wantLiteral, tok.Literal)
		}
	}
}

func TestTwoCharOperators(t *testing.T) {
	tests := []struct {
		input   string
		wantTyp token.TokenType
		wantLit string
	}{
		{"==", token.EQ, "=="},
		{"!=", token.NOT_EQ, "!="},
		{"<=", token.LTE, "<="},
		{">=", token.GTE, ">="},
		{"&&", token.AND, "&&"},
		{"||", token.OR, "||"},
	}
	for _, tt := range tests {
		tok := New(tt.input).NextToken()
		if tok.Type != tt.wantTyp || tok.Literal != tt.wantLit {
			t.Errorf("input=%q: want {%s %s}, got {%s %s}", tt.input, tt.wantTyp, tt.wantLit, tok.Type, tok.Literal)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	cases := []struct{ input, want string }{
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"quote\"inside"`, `quote"inside`},
		{`"back\\slash"`, `back\slash`},
	}
	for _, c := range cases {
		tok := New(c.input).NextToken()
		if tok.Type != token.STRING {
			t.Errorf("input=%q: want STRING, got %s", c.input, tok.Type)
		}
		if tok.Literal != c.want {
			t.Errorf("input=%q: want %q, got %q", c.input, c.want, tok.Literal)
		}
	}
}

func TestLineComments(t *testing.T) {
	input := `// this is a comment
let x = 1 // inline comment
x`
	l := New(input)
	toks := []token.TokenType{token.LET, token.IDENT, token.ASSIGN, token.INT, token.IDENT, token.EOF}
	for i, want := range toks {
		tok := l.NextToken()
		if tok.Type != want {
			t.Errorf("[%d] want %s, got %s", i, want, tok.Type)
		}
	}
}

func TestLineNumbers(t *testing.T) {
	input := "let x = 1\nlet y = 2"
	l := New(input)
	l.NextToken() // let  — line 1
	l.NextToken() // x
	l.NextToken() // =
	l.NextToken() // 1
	tok := l.NextToken() // let — line 2
	if tok.Line != 2 {
		t.Errorf("want line 2, got %d", tok.Line)
	}
}

func TestAllPunctuation(t *testing.T) {
	input := "(){}[],;:"
	want := []token.TokenType{
		token.LPAREN, token.RPAREN,
		token.LBRACE, token.RBRACE,
		token.LBRACKET, token.RBRACKET,
		token.COMMA, token.SEMICOLON, token.COLON,
		token.EOF,
	}
	l := New(input)
	for i, w := range want {
		tok := l.NextToken()
		if tok.Type != w {
			t.Errorf("[%d] want %s, got %s", i, w, tok.Type)
		}
	}
}

func TestArithmeticOperators(t *testing.T) {
	input := "+ - * / %"
	want := []token.TokenType{token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.PERCENT}
	l := New(input)
	for _, w := range want {
		tok := l.NextToken()
		if tok.Type != w {
			t.Errorf("want %s, got %s", w, tok.Type)
		}
	}
}
