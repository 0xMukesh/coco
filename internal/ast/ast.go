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
	}

	return out.String()
}

type Identifier struct {
	Token tokens.Token // tokens.IDENTIFIER
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token tokens.Token // tokens.INTEGER
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return fmt.Sprint(i.Value) }

type FloatLiteral struct {
	Token tokens.Token // tokens.FLOAT
	Value float64
}

func (f *FloatLiteral) expressionNode()      {}
func (f *FloatLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FloatLiteral) String() string       { return fmt.Sprint(f.Value) }

type UnaryExpression struct {
	Token      tokens.Token // tokens.BANG or tokens.MINUS
	Expression Expression
}

func (u *UnaryExpression) expressionNode()      {}
func (u *UnaryExpression) TokenLiteral() string { return u.Token.Literal }
func (u *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(u.Token.Literal)
	out.WriteString(" ")
	out.WriteString(u.Expression.String())
	out.WriteString(")")

	return out.String()
}

type LetStatement struct {
	Token      tokens.Token // tokens.LET
	Identifier *Identifier
	Value      Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Identifier.String())
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token      tokens.Token // tokens.RETURN
	Expression Expression
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")
	if s.Expression != nil {
		out.WriteString(s.Expression.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}
