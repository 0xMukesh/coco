package lexer

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()
	return &l
}

func (l *Lexer) readChar() byte {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
	return l.ch
}

func (l *Lexer) seekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readContinuous(f func(byte) bool) string {
	starting := l.position

	for f(l.ch) {
		l.readChar()
	}

	l.readPosition = l.position
	return l.input[starting:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() tokens.Token {
	var token tokens.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.seekChar() == '=' {
			ch := l.ch
			l.readChar()
			token = tokens.New(tokens.EQUALS, string(ch)+string(l.ch))
		} else {
			token = tokens.New(tokens.ASSIGN, string(l.ch))
		}
	case '+':
		token = tokens.New(tokens.PLUS, string(l.ch))
	case '-':
		token = tokens.New(tokens.MINUS, string(l.ch))
	case ',':
		token = tokens.New(tokens.COMMA, string(l.ch))
	case ';':
		token = tokens.New(tokens.SEMICOLON, string(l.ch))
	case '!':
		if l.seekChar() == '=' {
			ch := l.ch
			l.readChar()
			token = tokens.New(tokens.NOT_EQUALS, string(ch)+string(l.ch))
		} else {
			token = tokens.New(tokens.BANG, string(l.ch))
		}
	case '/':
		token = tokens.New(tokens.SLASH, string(l.ch))
	case '*':
		token = tokens.New(tokens.ASTERISK, string(l.ch))
	case '<':
		token = tokens.New(tokens.LESS_THAN, string(l.ch))
	case '>':
		token = tokens.New(tokens.GREATER_THAN, string(l.ch))
	case '(':
		token = tokens.New(tokens.LPAREN, string(l.ch))
	case ')':
		token = tokens.New(tokens.RPAREN, string(l.ch))
	case '{':
		token = tokens.New(tokens.LBRACE, string(l.ch))
	case '}':
		token = tokens.New(tokens.RBRACE, string(l.ch))
	case 0:
		token = tokens.New(tokens.EOF, "")
	default:
		if utils.IsLetter(l.ch) {
			identStr := l.readContinuous(utils.IsLetter)
			token = tokens.New(tokens.GetIdentTokenTypeByLiteral(identStr), identStr)
		} else if utils.IsNumber(l.ch) {
			num := l.readContinuous(utils.IsNumber)
			token = tokens.New(tokens.INT, num)
		} else {
			token = tokens.New(tokens.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return token
}

func (l *Lexer) String() string {
	return fmt.Sprintf("Lexer{input:%s, position:%d, readPosition:%d, ch:%q}", l.input, l.position, l.readPosition, l.ch)
}
