package ast

import (
	"bytes"
	"fmt"

	"github.com/0xmukesh/coco/internal/tokens"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}

	return out.String()
}

type IdentifierExpression struct {
	Token tokens.Token
	Value string
}

func (ie *IdentifierExpression) expressionNode() {}
func (ie *IdentifierExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IdentifierExpression) String() string {
	return ie.Value
}

type IntegerExpression struct {
	Token tokens.Token
	Value int64
}

func (ie *IntegerExpression) expressionNode() {}
func (ie *IntegerExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IntegerExpression) String() string {
	return fmt.Sprint(ie.Value)
}

type FloatExpression struct {
	Token tokens.Token
	Value float64
}

func (fe *FloatExpression) expressionNode() {}
func (fe *FloatExpression) TokenLiteral() string {
	return fe.Token.Literal
}
func (fe *FloatExpression) String() string {
	return fmt.Sprint(fe.Value)
}

// <prefix><expression>
type UnaryExpression struct {
	Token tokens.Token
	Expr  Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}
func (ue *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ue.TokenLiteral())
	out.WriteString(ue.Expr.String())

	return out.String()
}

// let <identifier> = <value>
// let <identifier>
type LetStatement struct {
	Token      tokens.Token
	Identifier *IdentifierExpression
	Value      Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Identifier.String())

	if ls.Value != nil {
		out.WriteString(" = ")
		out.WriteString(ls.Value.String())
	}

	return out.String()
}

// return <expr>
type ReturnStatement struct {
	Token tokens.Token
	Expr  Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())

	if rs.Expr != nil {
		out.WriteString(" ")
		out.WriteString(rs.Expr.String())
	}

	return out.String()
}

// wrapper statement type for expression, so that expressions can be executed directly via REPL
type ExpressionStatement struct {
	Token tokens.Token
	Expr  Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	return es.Expr.String()
}
