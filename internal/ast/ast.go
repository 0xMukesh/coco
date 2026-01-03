package ast

import (
	"bytes"
	"fmt"
	"strings"

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

type BooleanLiteral struct {
	Token tokens.Token // tokens.BOOLEAN
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return fmt.Sprint(b.Value) }

type UnaryExpression struct {
	Token      tokens.Token // tokens.BANG or tokens.MINUS
	Expression Expression
}

func (u *UnaryExpression) expressionNode()      {}
func (u *UnaryExpression) TokenLiteral() string { return u.Token.Literal }
func (u *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(u.Token.Literal + " ")
	out.WriteString(u.Expression.String())
	out.WriteString(")")

	return out.String()
}

type BinaryExpression struct {
	Token tokens.Token
	Left  Expression
	Right Expression
}

func (b *BinaryExpression) expressionNode()      {}
func (b *BinaryExpression) TokenLiteral() string { return b.Token.Literal }
func (b *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(b.Left.String() + " ")
	out.WriteString(b.Token.Literal + " ")
	out.WriteString(b.Right.String())
	out.WriteString(")")

	return out.String()
}

// in coco, the results of a conditional expression can be binded to be a variable
type IfExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer

	// if <condition> <consquence> (else <alternative>)
	out.WriteString("if")
	out.WriteString(" " + i.Condition.String())
	out.WriteString(" " + i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString(" else")
		out.WriteString(" " + i.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      tokens.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     tokens.Token // tokens.LPAREN
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
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

	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token // tokens.LBRACE; start of a block
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range b.Statements {
		out.WriteString(stmt.String())
	}

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
