package lexer

import (
	"bytes"
	"errors"
	"strings"

	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
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

	l.readChar()
	return l
}

func (l *Lexer) newToken(tt tokens.TokenType, literal string) tokens.Token {
	endColumn := l.column + len(literal)

	if tt == tokens.ILLEGAL {
		endColumn = l.column + 1
	}

	return tokens.New(tt, literal, l.line, l.column, endColumn)
}

func (l *Lexer) newTokenWithExplicitStartColumn(tt tokens.TokenType, startColumn int, literal string) tokens.Token {
	return tokens.New(tt, literal, l.line, startColumn, startColumn+len(literal))
}

func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' || l.currChar == '\n' || l.currChar == '\t' || l.currChar == '\r' {
		l.readChar()
	}
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

func (l *Lexer) readIdentifier() string {
	startPosition := l.currPosition

	// check if the next character is letter or number
	// if yes, then consume it
	for utils.IsLetter(l.peekChar()) || utils.IsDigit(l.peekChar()) {
		l.readChar()
	}

	return l.input[startPosition : l.currPosition+1]
}

func (l *Lexer) readNumeric(skipDot bool) string {
	startPosition := l.currPosition
	decimalSeparatorSeen := false

	// check if next character is numeric i.e. is a either a digit or "."
	// if yes, then consume it
	for utils.IsDigit(l.peekChar()) || (!skipDot && !decimalSeparatorSeen && l.peekChar() == '.') {
		if l.peekChar() == '.' {
			decimalSeparatorSeen = true
		}

		l.readChar()
	}

	return l.input[startPosition : l.currPosition+1]
}

func (l *Lexer) readString(delim byte) (string, error) {
	var out bytes.Buffer
	foundClosingDelim := false

	out.WriteByte(delim) // writing starting delim to the token literal

	for {
		l.readChar()

		if l.currChar == delim {
			foundClosingDelim = true
		}

		if l.currChar == delim || l.currChar == 0 {
			out.WriteByte(l.currChar) // writing closing delim to the token literal
			break
		}

		if l.currChar == '\\' {
			if utils.IsEscapeSequenceCode(l.peekChar()) {
				escapeSequenceCode := l.peekChar()
				l.readChar()

				if escapeSequenceCode == 'n' {
					out.WriteByte('\n')
				} else if escapeSequenceCode == 't' {
					out.WriteByte('\t')
				} else if escapeSequenceCode == '"' && delim == '"' {
					out.WriteByte('"')
				}

				continue
			} else {
				return "", errors.New("invalid escape character")
			}
		}

		out.WriteByte(l.currChar)
	}

	if foundClosingDelim {
		return out.String(), nil
	} else {
		return out.String(), errors.New("unterminated string")
	}
}

func (l *Lexer) nextToken() tokens.Token {
	var tok tokens.Token

	l.skipWhitespace()

	switch l.currChar {
	case '+':
		startColumn := l.column

		if l.peekChar() == '+' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.INCREMENT, startColumn, "++")
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.PLUS_EQUAL, startColumn, "+=")
		} else {
			tok = l.newToken(tokens.PLUS, string(l.currChar))
		}
	case '-':
		startColumn := l.column

		if l.peekChar() == '-' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.DECREMENT, startColumn, "--")
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.MINUS_EQUAL, startColumn, "-=")
		} else {
			tok = l.newToken(tokens.MINUS, string(l.currChar))
		}
	case '*':
		startColumn := l.column

		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.STAR_EQUAL, startColumn, "*=")
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.DOUBLE_STAR, startColumn, "**")
		} else {
			tok = l.newToken(tokens.STAR, string(l.currChar))
		}
	case '/':
		if l.peekChar() == '/' {
			// current character is slash
			l.readChar()

			for l.currChar != '\n' && l.currChar != 0 {
				// skips all the characters until new line
				l.readChar()
			}

			// once new line is found `NextToken` is executed recursively, which runs `skipWhitespace`, which under the hood gets the next character
			return l.nextToken()
		} else if l.peekChar() == '*' {
			// current character is star
			l.readChar()
			// consume star
			l.readChar()

			for {
				if l.currChar == 0 {
					return l.newToken(tokens.ILLEGAL, "unterminated comment")
				}

				if l.currChar == '*' && l.peekChar() == '/' {
					// consume closing star
					l.readChar()
					// consume closing slash
					l.readChar()

					return l.nextToken()
				}

				l.readChar()
			}
		} else if l.peekChar() == '=' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.SLASH_EQUAL, startColumn, "/=")
		} else {
			tok = l.newToken(tokens.SLASH, string(l.currChar))
		}
	case '%':
		tok = l.newToken(tokens.MODULO, "%")
	case '=':
		if l.peekChar() == '=' {
			startColumn := l.column
			// consume assign token
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.EQUALS, startColumn, "==")
		} else {
			tok = l.newToken(tokens.ASSIGN, string(l.currChar))
		}
	case '<':
		if l.peekChar() == '=' {
			startColumn := l.column
			// consume lt token
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.LESS_THAN_EQUALS, startColumn, "<=")
		} else {
			tok = l.newToken(tokens.LESS_THAN, string(l.currChar))
		}
	case '>':
		if l.peekChar() == '=' {
			startColumn := l.column
			// consume gt token
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.GREATER_THAN_EQUALS, startColumn, ">=")
		} else {
			tok = l.newToken(tokens.GREATER_THAN, string(l.currChar))
		}
	case '!':
		if l.peekChar() == '=' {
			startColumn := l.column
			// consume bang token
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.NOT_EQUALS, startColumn, "!=")
		} else {
			tok = l.newToken(tokens.BANG, string(l.currChar))
		}
	case '&':
		if l.peekChar() == '&' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.AND, startColumn, "&&")
		}
	case '|':
		if l.peekChar() == '|' {
			startColumn := l.column
			l.readChar()
			tok = l.newTokenWithExplicitStartColumn(tokens.OR, startColumn, "||")
		}
	case '(':
		tok = l.newToken(tokens.LPAREN, string(l.currChar))
	case ')':
		tok = l.newToken(tokens.RPAREN, string(l.currChar))
	case '{':
		tok = l.newToken(tokens.LBRACE, string(l.currChar))
	case '}':
		tok = l.newToken(tokens.RBRACE, string(l.currChar))
	case '[':
		tok = l.newToken(tokens.LSQUARE, string(l.currChar))
	case ']':
		tok = l.newToken(tokens.RSQUARE, string(l.currChar))
	case ';':
		tok = l.newToken(tokens.SEMICOLON, string(l.currChar))
	case ',':
		tok = l.newToken(tokens.COMMA, string(l.currChar))
	case ':':
		tok = l.newToken(tokens.COLON, string(l.currChar))
	case '"':
		startColumn := l.column + 1
		str, err := l.readString(l.currChar)

		if err != nil {
			tok = l.newTokenWithExplicitStartColumn(tokens.ILLEGAL, startColumn, err.Error())
		} else {
			tok = l.newTokenWithExplicitStartColumn(tokens.STRING, startColumn, str)
		}
	case 0:
		tok = l.newToken(tokens.EOF, "")
	default:
		if utils.IsLetter(l.currChar) {
			startColumn := l.column
			identifier := l.readIdentifier()

			tok = l.newTokenWithExplicitStartColumn(tokens.IdentTokenTypeLookup(identifier), startColumn, identifier)
		} else if utils.IsDigit(l.currChar) {
			startColumn := l.column
			numeric := l.readNumeric(false)

			if strings.Contains(numeric, ".") {
				tok = l.newTokenWithExplicitStartColumn(tokens.FLOAT, startColumn, numeric)
			} else {
				tok = l.newTokenWithExplicitStartColumn(tokens.INTEGER, startColumn, numeric)
			}
		} else {
			tok = l.newToken(tokens.ILLEGAL, string(l.currChar))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) Lex() []tokens.Token {
	var tks []tokens.Token

	for tok := l.nextToken(); tok.Type != tokens.EOF; tok = l.nextToken() {
		tks = append(tks, tok)
	}

	return tks
}
