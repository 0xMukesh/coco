package tokens

type TokenType string

type Token struct {
	Type          TokenType
	Literal       string
	Line          int
	StartPosition int
	EndPosition   int
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
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
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

func New(tokenType TokenType, literal string, line, startPosition, endPosition int) Token {
	return Token{
		Type:          tokenType,
		Literal:       literal,
		Line:          line,
		StartPosition: startPosition,
		EndPosition:   endPosition,
	}
}

func GetIdentTokenTypeByLiteral(literal string) TokenType {
	if tok, ok := KEYWORDS[literal]; ok {
		return tok
	}

	return IDENTIFIER
}
