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
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"

	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"
	ASSIGN = "="

	COMMA     = ","
	SEMICOLON = ";"
	QUOTES    = "\""

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
)

var keywords = map[string]TokenType{
	"fn":    FUNCTION,
	"let":   LET,
	"const": CONST,
}

func LookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}

	return IDENTIFIER
}

func New(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}
