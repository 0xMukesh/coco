package lexer

import (
	"github.com/0xmukesh/coco/internal/tokens"
)

type Lexer struct {
	input        string
	currPosition int
	nextPosition int
	currChar     byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		line:         1,
		currPosition: 0,
		column:       -1,
	}

	l.readChar() // initialize with currChar
	return l
}

func (l *Lexer) newToken(tt tokens.TokenType, literal string) tokens.Token {
	return tokens.New(tt, literal, l.line, l.column, l.column+len(literal))
}

func (l *Lexer) newTokenWithExplicitStartColumn(tt tokens.TokenType, startColumn int, literal string) tokens.Token {
	return tokens.New(tt, literal, l.line, startColumn, startColumn+len(literal))
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.currChar = 0
	} else {
		l.currChar = l.input[l.nextPosition]

		if l.currChar == '\n' {
			l.line++
			l.column = -1
		} else {
			l.column++
		}
	}

	l.currPosition = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' || l.currChar == '\n' || l.currChar == '\t' || l.currChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() tokens.Token {
	var tok tokens.Token

	l.skipWhitespace()

	switch l.currChar {
	case '+':
		tok = l.newToken(tokens.PLUS, string(l.currChar))
	case '-':
		tok = l.newToken(tokens.MINUS, string(l.currChar))
	case '*':
		tok = l.newToken(tokens.STAR, string(l.currChar))
	case '/':
		tok = l.newToken(tokens.SLASH, string(l.currChar))
	case '=':
		if l.peekChar() == '=' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.EQUALS, startColumn, "==")
		} else {
			tok = l.newToken(tokens.ASSIGN, string(l.currChar))
		}
	case '<':
		if l.peekChar() == '=' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.LESS_THAN_EQUALS, startColumn, "<=")
		} else {
			tok = l.newToken(tokens.LESS_THAN, string(l.currChar))
		}
	case '>':
		if l.peekChar() == '=' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.GREATER_THAN_EQUALS, startColumn, ">=")
		} else {
			tok = l.newToken(tokens.GREATER_THAN, string(l.currChar))
		}
	case '!':
		if l.peekChar() == '=' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.NOT_EQUALS, startColumn, "!=")
		} else {
			tok = l.newToken(tokens.BANG, string(l.currChar))
		}
	case '(':
		tok = l.newToken(tokens.LPAREN, string(l.currChar))
	case ')':
		tok = l.newToken(tokens.RPAREN, string(l.currChar))
	case '{':
		tok = l.newToken(tokens.LBRACE, string(l.currChar))
	case '}':
		tok = l.newToken(tokens.RBRACE, string(l.currChar))
	case 0:
		tok = l.newToken(tokens.EOF, "")
	default:
		tok = l.newToken(tokens.ILLEGAL, string(l.currChar))
	}

	l.readChar()
	return tok
}
