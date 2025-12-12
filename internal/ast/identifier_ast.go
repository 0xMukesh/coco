package ast

import "github.com/0xmukesh/coco/internal/tokens"

type IdentifierExpression struct {
	Token tokens.Token
	Value string
}

func (i *IdentifierExpression) expressionNode() {}
func (i *IdentifierExpression) NodeLiteral() string {
	return i.Token.Literal
}
