package tokens

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"

	PLUS               = "+"
	MINUS              = "-"
	STAR               = "*"
	SLASH              = "/"
	ASSIGN             = "="
	LESS_THAN          = "<"
	LESS_THAN_EQUAL    = "<="
	GREATER_THAN       = ">"
	GREATER_THAN_EQUAL = ">="
	BANG               = "!"
	EQUALS             = "=="
	NOT_EQUALS         = "!="

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
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"const":  CONST,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}

	return IDENTIFIER
}

func New(tokenType TokenType, literal string, line int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}
