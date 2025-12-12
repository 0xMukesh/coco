package tokens

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	FLOAT      = "FLOAT"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"

	COMMA     = ","
	SEMICOLON = ";"
	BANG      = "!"
	SLASH     = "/"
	ASTERISK  = "*"

	LESS_THAN    = "<"
	GREATER_THAN = ">"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	EQUALS     = "=="
	NOT_EQUALS = "!="

	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONSTANT = "CONSTANT"
	TRUE     = "TRUE"
	FALSE    = "false"
	IF       = "if"
	ELSE     = "else"
	RETURN   = "return"
)

var KEYWORDS = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"const":  CONSTANT,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func New(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func GetIdentTokenTypeByLiteral(literal string) TokenType {
	if tok, ok := KEYWORDS[literal]; ok {
		return tok
	}

	return IDENTIFIER
}
