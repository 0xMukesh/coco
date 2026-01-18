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
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"
	MODULO = "%"

	LESS_THAN           = "<"
	GREATER_THAN        = ">"
	EQUALS              = "=="
	LESS_THAN_EQUALS    = "<="
	GREATER_THAN_EQUALS = ">="
	NOT_EQUALS          = "!="
	AND                 = "&&"
	OR                  = "||"

	ASSIGN = "="
	BANG   = "!"

	INCREMENT   = "++"
	DECREMENT   = "--"
	DOUBLE_STAR = "**"
	PLUS_EQUAL  = "+="
	MINUS_EQUAL = "-="
	STAR_EQUAL  = "*="
	SLASH_EQUAL = "/="

	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	LSQUARE = "["
	RSQUARE = "]"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"
	STRING     = "STRING"

	LET      = "LET"
	CONST    = "CONST"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	EXIT     = "EXIT"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var PROGRAM_KEYWORDS = map[string]TokenType{
	"let":      LET,
	"const":    CONST,
	"fn":       FUNCTION,
	"if":       IF,
	"else":     ELSE,
	"true":     TRUE,
	"false":    FALSE,
	"return":   RETURN,
	"for":      FOR,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
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

func NewMinimal(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func IdentTokenTypeLookup(ident string) TokenType {
	if tt, ok := PROGRAM_KEYWORDS[ident]; ok {
		return tt
	}

	return IDENTIFIER
}
