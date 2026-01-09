package tokens

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
	Line    int
	Column  int
}

const (
	PLUS  = "+"
	MINUS = "-"
	STAR  = "*"
	SLASH = "/"

	ASSIGN              = "="
	LESS_THAN           = "<"
	GREATER_THAN        = ">"
	EQUALS              = "=="
	LESS_THAN_EQUALS    = "<="
	GREATER_THAN_EQUALS = ">="
	BANG                = "!"
	NOT_EQUALS          = "!="

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	SEMICOLON = ";"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

func New(tokenType TokenType, literal string, line, column int) Token {
	return Token{
		Literal: literal,
		Type:    tokenType,
		Line:    line,
		Column:  column,
	}
}
