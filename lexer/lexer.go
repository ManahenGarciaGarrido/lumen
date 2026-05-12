package lexer

import "github.com/ManahenGarciaGarrido/lumen/token"

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int
	col     int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	// skip line comments
	for l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
	}

	ln, col := l.line, l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.EQ, Literal: "==", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.ASSIGN, Literal: "=", Line: ln, Col: col}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.NOT_EQ, Literal: "!=", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.BANG, Literal: "!", Line: ln, Col: col}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.LTE, Literal: "<=", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.LT, Literal: "<", Line: ln, Col: col}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.GTE, Literal: ">=", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.GT, Literal: ">", Line: ln, Col: col}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.AND, Literal: "&&", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.ILLEGAL, Literal: "&", Line: ln, Col: col}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			l.readChar()
			return token.Token{Type: token.OR, Literal: "||", Line: ln, Col: col}
		}
		l.readChar()
		return token.Token{Type: token.ILLEGAL, Literal: "|", Line: ln, Col: col}
	case '+':
		l.readChar()
		return token.Token{Type: token.PLUS, Literal: "+", Line: ln, Col: col}
	case '-':
		l.readChar()
		return token.Token{Type: token.MINUS, Literal: "-", Line: ln, Col: col}
	case '*':
		l.readChar()
		return token.Token{Type: token.ASTERISK, Literal: "*", Line: ln, Col: col}
	case '/':
		l.readChar()
		return token.Token{Type: token.SLASH, Literal: "/", Line: ln, Col: col}
	case '%':
		l.readChar()
		return token.Token{Type: token.PERCENT, Literal: "%", Line: ln, Col: col}
	case ',':
		l.readChar()
		return token.Token{Type: token.COMMA, Literal: ",", Line: ln, Col: col}
	case ';':
		l.readChar()
		return token.Token{Type: token.SEMICOLON, Literal: ";", Line: ln, Col: col}
	case ':':
		l.readChar()
		return token.Token{Type: token.COLON, Literal: ":", Line: ln, Col: col}
	case '(':
		l.readChar()
		return token.Token{Type: token.LPAREN, Literal: "(", Line: ln, Col: col}
	case ')':
		l.readChar()
		return token.Token{Type: token.RPAREN, Literal: ")", Line: ln, Col: col}
	case '{':
		l.readChar()
		return token.Token{Type: token.LBRACE, Literal: "{", Line: ln, Col: col}
	case '}':
		l.readChar()
		return token.Token{Type: token.RBRACE, Literal: "}", Line: ln, Col: col}
	case '[':
		l.readChar()
		return token.Token{Type: token.LBRACKET, Literal: "[", Line: ln, Col: col}
	case ']':
		l.readChar()
		return token.Token{Type: token.RBRACKET, Literal: "]", Line: ln, Col: col}
	case '"':
		str, ok := l.readString()
		if !ok {
			return token.Token{Type: token.ILLEGAL, Literal: str, Line: ln, Col: col}
		}
		return token.Token{Type: token.STRING, Literal: str, Line: ln, Col: col}
	case 0:
		return token.Token{Type: token.EOF, Literal: "", Line: ln, Col: col}
	default:
		if isLetter(l.ch) {
			start, startLn, startCol := l.pos, l.line, l.col
			for isLetter(l.ch) || isDigit(l.ch) {
				l.readChar()
			}
			lit := l.input[start:l.pos]
			return token.Token{Type: token.LookupIdent(lit), Literal: lit, Line: startLn, Col: startCol}
		}
		if isDigit(l.ch) {
			start, startLn, startCol := l.pos, l.line, l.col
			for isDigit(l.ch) {
				l.readChar()
			}
			return token.Token{Type: token.INT, Literal: l.input[start:l.pos], Line: startLn, Col: startCol}
		}
		ch := l.ch
		l.readChar()
		return token.Token{Type: token.ILLEGAL, Literal: string(ch), Line: ln, Col: col}
	}
}

func (l *Lexer) readString() (string, bool) {
	l.readChar() // skip opening "
	var out []byte
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out = append(out, '\n')
			case 't':
				out = append(out, '\t')
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			default:
				out = append(out, '\\', l.ch)
			}
		} else {
			out = append(out, l.ch)
		}
		l.readChar()
	}
	if l.ch == 0 {
		return string(out), false
	}
	l.readChar() // skip closing "
	return string(out), true
}

func isLetter(ch byte) bool { return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_' }
func isDigit(ch byte) bool  { return ch >= '0' && ch <= '9' }
