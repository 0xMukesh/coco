package lexer

import (
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type Lexer struct {
	input           string
	currentPosition int
	peekPosition    int
	currChar        byte
}

func New(input string) *Lexer {
	l := Lexer{input: input, currentPosition: 0, peekPosition: 0}
	l.readChar()

	return &l
}

func (l *Lexer) Tokenize() []tokens.Token {
	tks := []tokens.Token{}

	for {
		tok := l.nextToken()
		tks = append(tks, tok)

		if tok.Type == tokens.EOF {
			break
		}
	}

	return tks
}

func (l *Lexer) readChar() byte {
	if l.peekPosition >= len(l.input) {
		l.currChar = 0
	} else {
		l.currChar = l.input[l.peekPosition]
	}

	l.currentPosition = l.peekPosition
	l.peekPosition++

	return l.currChar
}

func (l *Lexer) seekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}

func (l *Lexer) readContinuous(f func(byte) bool) string {
	starting := l.currentPosition

	for f(l.seekChar()) {
		l.readChar()
	}

	return l.input[starting:l.peekPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' || l.currChar == '\t' || l.currChar == '\n' || l.currChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) nextToken() tokens.Token {
	var token tokens.Token

	l.skipWhitespace()

	switch l.currChar {
	case '=':
		if l.seekChar() == '=' {
			l.readChar()
			token = l.constructMultiCharToken("==", tokens.EQUALS)
		} else {
			token = l.constructSingleCharToken(tokens.ASSIGN)
		}
	case '+':
		token = l.constructSingleCharToken(tokens.PLUS)
	case '-':
		token = l.constructSingleCharToken(tokens.MINUS)
	case ',':
		token = l.constructSingleCharToken(tokens.COMMA)
	case ';':
		token = l.constructSingleCharToken(tokens.SEMICOLON)
	case '!':
		if l.seekChar() == '=' {
			l.readChar()
			token = l.constructMultiCharToken("!=", tokens.NOT_EQUALS)
		} else {
			token = l.constructSingleCharToken(tokens.BANG)
		}
	case '/':
		token = l.constructSingleCharToken(tokens.SLASH)
	case '*':
		token = l.constructSingleCharToken(tokens.ASTERISK)
	case '<':
		token = l.constructSingleCharToken(tokens.LESS_THAN)
	case '>':
		token = l.constructSingleCharToken(tokens.GREATER_THAN)
	case '(':
		token = l.constructSingleCharToken(tokens.LPAREN)
	case ')':
		token = l.constructSingleCharToken(tokens.RPAREN)
	case '{':
		token = l.constructSingleCharToken(tokens.LBRACE)
	case '}':
		token = l.constructSingleCharToken(tokens.RBRACE)
	case 0:
		token = l.constructSingleCharToken(tokens.EOF)
	default:
		if utils.IsLetter(l.currChar) {
			identStr := l.readContinuous(utils.IsLetter)
			token = l.constructMultiCharToken(identStr, tokens.GetIdentTokenTypeByLiteral(identStr))
		} else if utils.IsNumber(l.currChar) {
			num := l.readContinuous(utils.IsNumber)
			token = l.constructMultiCharToken(num, tokens.INT)
		} else {
			token = l.constructSingleCharToken(tokens.ILLEGAL)
		}
	}

	l.readChar()
	return token
}

func (l *Lexer) constructSingleCharToken(tokenType tokens.TokenType) tokens.Token {
	return tokens.New(tokenType, string(l.currChar))
}

func (l *Lexer) constructMultiCharToken(literal string, tokenType tokens.TokenType) tokens.Token {
	return tokens.New(tokenType, literal)
}
