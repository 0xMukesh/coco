package ast

import "github.com/0xmukesh/coco/internal/tokens"

type Node interface {
	NodeLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// `let x = 5;` is a statement
// `5;` is a expression and the value produced is 5
type Program struct {
	Statements  []Statement  // statements don't produce values
	Expressions []Expression // expression produce values
}

func (p Program) NodeLiteral() string {
	return "RootNode"
}

// expressions
type IdentifierExpression struct {
	Token tokens.Token
	Value string
}

func (i *IdentifierExpression) expressionNode() {}
func (i *IdentifierExpression) NodeLiteral() string {
	return i.Token.Literal
}

// statements
type LetStatement struct {
	Token      tokens.Token
	Identifier *IdentifierExpression
	Value      Expression
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) NodeLiteral() string {
	return l.Token.Literal
}

type ReturnStatement struct {
	Token tokens.Token
	Value Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) NodeLiteral() string {
	return r.Token.Literal
}
