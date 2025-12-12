package ast

import (
	"github.com/0xmukesh/coco/internal/tokens"
)

type LetStatement struct {
	Token      tokens.Token
	Identifier *IdentifierExpression
	Value      Expression
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) NodeLiteral() string {
	return l.Token.Literal
}
