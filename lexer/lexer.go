package lexer

import (
	"bufio"
	"fmt"
	"monkey/token"
	"unicode"
)

type Lexer struct {
	Input        string
	position     int
	nextPosition int
	char         byte
	scanner      *bufio.Scanner
}

func New(Input string, scanner *bufio.Scanner) *Lexer {
	lexer := &Lexer{Input: Input, scanner: scanner}
	lexer.ReadChar()
	return lexer
}

func (l *Lexer) ReadChar() {

	if l.nextPosition >= len(l.Input) {
		l.char = 0

		if l.scanner.Scan() {
			if len(l.scanner.Text()) == 0 {
				l.ReadChar()
				return
			}
			l.Input = l.Input + " " + l.scanner.Text()
			l.char = l.Input[l.nextPosition]
		}

	} else {
		l.char = l.Input[l.nextPosition]
	}

	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) lookAhead() byte {
	if l.nextPosition > len(l.Input) {
		return 0
	}
	return l.Input[l.nextPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.char {
	case '=':
		if l.lookAhead() == '=' {
			tok = token.Token{Type: token.EQUAL, Literal: string(l.char) + string(l.Input[l.nextPosition])}
			l.ReadChar()
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case '!':
		if l.lookAhead() == '=' {
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(l.char) + string(l.Input[l.nextPosition])}
			l.ReadChar()
		} else {
			tok = newToken(token.NOT, l.char)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '<':
		tok = newToken(token.LESS_THAN, l.char)
	case '>':
		tok = newToken(token.GREATER_THAN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '[':
		tok = newToken(token.LSQBRACKET, l.char)
	case ']':
		tok = newToken(token.RSQBRACKET, l.char)
	case '.':
		tok = newToken(token.DOT, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '/':
		tok = newToken(token.DIVIDE, l.char)
	case '*':
		tok = newToken(token.MULTIPLY, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if unicode.IsLetter(rune(l.char)) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.GetIdentType(tok.Literal)
			return tok
		} else if unicode.IsNumber(rune(l.char)) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.ReadChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(rune(l.char)) {
		l.ReadChar()
	}
	return l.Input[position:l.position]
}

func (l *Lexer) readString() string {
	l.ReadChar()
	position := l.position

	for l.char != '"' {
		if l.char == 0 {
			fmt.Println("No string ending character")
			return ""
		}
		l.ReadChar()
	}
	return l.Input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsNumber(rune(l.char)) {
		l.ReadChar()
	}
	return l.Input[position:l.position]
}

func newToken(tokenType string, tokenValue byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(tokenValue)}
}

func (l *Lexer) eatWhitespace() {
	for l.char == ' ' || l.char == '\n' || l.char == '\r' || l.char == '\t' {
		l.ReadChar()
	}
}
