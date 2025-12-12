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

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONSTANT = "CONSTANT"
)

var KEYWORDS = map[string]TokenType{
	"fn":    FUNCTION,
	"let":   LET,
	"const": CONSTANT,
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
