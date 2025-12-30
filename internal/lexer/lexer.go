package lexer

import (
	"strings"

	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type Lexer struct {
	input           string
	currentPosition int
	peekPosition    int
	currChar        byte
	currLine        int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:           input,
		currLine:        1,
		currentPosition: 0,
	}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.peekPosition >= len(l.input) {
		l.currChar = 0
	} else {
		l.currChar = l.input[l.peekPosition]
	}

	l.currentPosition = l.peekPosition
	l.peekPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.currentPosition
	for utils.IsLetter(l.currChar) {
		l.readChar()
	}

	return l.input[startPosition:l.currentPosition]
}

func (l *Lexer) readDigit() string {
	startPosition := l.currentPosition
	for utils.IsDigit(l.currChar) {
		l.readChar()
	}

	return l.input[startPosition:l.currentPosition]
}

func (l *Lexer) skipWhitespace() {
	for utils.IsWhitespace(l.currChar) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() tokens.Token {
	var tok tokens.Token

	if l.currChar == '\n' {
		l.currLine++
	}

	l.skipWhitespace()

	switch l.currChar {
	case '+':
		tok = tokens.New(tokens.PLUS, string(l.currChar), l.currLine)
	case '-':
		tok = tokens.New(tokens.MINUS, string(l.currChar), l.currLine)
	case '*':
		tok = tokens.New(tokens.STAR, string(l.currChar), l.currLine)
	case '/':
		tok = tokens.New(tokens.SLASH, string(l.currChar), l.currLine)
	case '=':
		if l.peekChar() == '=' {
			tok = tokens.New(tokens.EQUALS, "==", l.currLine)
			l.readChar()
		} else {
			tok = tokens.New(tokens.ASSIGN, string(l.currChar), l.currLine)
		}
	case '<':
		tok = tokens.New(tokens.LESS_THAN, string(l.currChar), l.currLine)
	case '>':
		tok = tokens.New(tokens.GREATER_THAN, string(l.currChar), l.currLine)
	case '!':
		if l.peekChar() == '=' {
			tok = tokens.New(tokens.NOT_EQUALS, "!=", l.currLine)
			l.readChar()
		} else {
			tok = tokens.New(tokens.BANG, string(l.currChar), l.currLine)
		}
	case ',':
		tok = tokens.New(tokens.COMMA, string(l.currChar), l.currLine)
	case ';':
		tok = tokens.New(tokens.SEMICOLON, string(l.currChar), l.currLine)
	case '"':
		tok = tokens.New(tokens.QUOTES, string(l.currChar), l.currLine)
	case '(':
		tok = tokens.New(tokens.LPAREN, string(l.currChar), l.currLine)
	case ')':
		tok = tokens.New(tokens.RPAREN, string(l.currChar), l.currLine)
	case '{':
		tok = tokens.New(tokens.LBRACE, string(l.currChar), l.currLine)
	case '}':
		tok = tokens.New(tokens.RBRACE, string(l.currChar), l.currLine)
	case 0:
		tok = tokens.New(tokens.EOF, string(l.currChar), l.currLine)
	default:
		if utils.IsLetter(l.currChar) {
			s := l.readIdentifier()
			tok = tokens.New(tokens.LookupIdent(s), s, l.currLine)
			// skip l.readChar call, as it happens at the last cycle of the loop
			return tok
		} else if utils.IsDigit(l.currChar) {
			d := l.readDigit()

			if strings.Contains(d, ".") {
				tok = tokens.New(tokens.FLOAT, d, l.currLine)
			} else {
				tok = tokens.New(tokens.INTEGER, d, l.currLine)
			}

			return tok
		} else {
			tok = tokens.New(tokens.ILLEGAL, string(l.currChar), l.currLine)
		}
	}

	l.readChar()
	return tok
}
