package tokens

type TokenType string

type Token struct {
	Literal     string
	Type        TokenType
	Line        int
	StartColumn int
	EndColumn   int
}

const (
	ASSIGN = "="
	BANG   = "!"

	PLUS  = "+"
	MINUS = "-"
	STAR  = "*"
	SLASH = "/"

	LESS_THAN           = "<"
	GREATER_THAN        = ">"
	EQUALS              = "=="
	LESS_THAN_EQUALS    = "<="
	GREATER_THAN_EQUALS = ">="
	NOT_EQUALS          = "!="

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	SEMICOLON = ";"
	COMMA     = ","

	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"

	LET      = "LET"
	CONST    = "CONST"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"const":  CONST,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func New(tokenType TokenType, literal string, line, startColumn, endColumn int) Token {
	return Token{
		Literal:     literal,
		Type:        tokenType,
		Line:        line,
		StartColumn: startColumn,
		EndColumn:   endColumn,
	}
}

func IdentTokenTypeLookup(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}

	return IDENTIFIER
}
