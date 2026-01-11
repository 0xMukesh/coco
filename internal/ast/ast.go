package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/types"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
	GetType() types.Type
	SetType(types.Type) types.Type
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
	Token   tokens.Token
	Literal string
	Type    types.Type
}

func (ie *IdentifierExpression) expressionNode() {}
func (ie *IdentifierExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IdentifierExpression) String() string {
	return ie.Literal
}
func (ie *IdentifierExpression) GetType() types.Type {
	return ie.Type
}
func (ie *IdentifierExpression) SetType(t types.Type) types.Type {
	ie.Type = t
	return t
}

type StringExpression struct {
	Token tokens.Token
	Value string
	Type  types.Type
}

func (se *StringExpression) expressionNode() {}
func (se *StringExpression) TokenLiteral() string {
	return se.Token.Literal
}
func (se *StringExpression) String() string {
	return se.Value
}
func (se *StringExpression) GetType() types.Type {
	return se.Type
}
func (se *StringExpression) SetType(t types.Type) types.Type {
	se.Type = t
	return t
}

type BooleanExpression struct {
	Token tokens.Token
	Value bool
	Type  types.Type
}

func (be *BooleanExpression) expressionNode() {}
func (be *BooleanExpression) TokenLiteral() string {
	return be.Token.Literal
}
func (be *BooleanExpression) String() string {
	return fmt.Sprint(be.Value)
}
func (be *BooleanExpression) GetType() types.Type {
	return be.Type
}
func (be *BooleanExpression) SetType(t types.Type) types.Type {
	be.Type = t
	return t
}

type IntegerExpression struct {
	Token tokens.Token
	Value int64
	Type  types.Type
}

func (ie *IntegerExpression) expressionNode() {}
func (ie *IntegerExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IntegerExpression) String() string {
	return fmt.Sprint(ie.Value)
}
func (ie *IntegerExpression) GetType() types.Type {
	return ie.Type
}
func (ie *IntegerExpression) SetType(t types.Type) types.Type {
	ie.Type = t
	return t
}

type FloatExpression struct {
	Token tokens.Token
	Value float64
	Type  types.Type
}

func (fe *FloatExpression) expressionNode() {}
func (fe *FloatExpression) TokenLiteral() string {
	return fe.Token.Literal
}
func (fe *FloatExpression) String() string {
	return fmt.Sprint(fe.Value)
}
func (fe *FloatExpression) GetType() types.Type {
	return fe.Type
}
func (fe *FloatExpression) SetType(t types.Type) types.Type {
	fe.Type = t
	return t
}

// <prefix><expression>
type UnaryExpression struct {
	Token tokens.Token
	Expr  Expression
	Type  types.Type
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
func (ue *UnaryExpression) GetType() types.Type {
	return ue.Type
}
func (ue *UnaryExpression) SetType(t types.Type) types.Type {
	ue.Type = t
	return t
}

// <left><operator><right>
type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
	Type     types.Type
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Operator.Literal
}
func (be *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.TokenLiteral() + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}
func (be *BinaryExpression) GetType() types.Type {
	return be.Type
}
func (be *BinaryExpression) SetType(t types.Type) types.Type {
	be.Type = t
	return t
}

type GroupedExpression struct {
	Token tokens.Token
	Expr  Expression
}

func (ge *GroupedExpression) expressionNode() {}
func (ge *GroupedExpression) TokenLiteral() string {
	return ge.Token.Literal
}
func (ge *GroupedExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ge.Expr.String())
	out.WriteString(")")

	return out.String()
}
func (ge *GroupedExpression) GetType() types.Type {
	return ge.Expr.GetType()
}
func (ge *GroupedExpression) SetType(t types.Type) types.Type {
	ge.Expr.SetType(t)
	return t
}

// if (<condition>) { <consequence> } ?(else { <alternative> })
// ?(...) = optional
type IfExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Type        types.Type
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
func (ie *IfExpression) GetType() types.Type {
	return ie.Type
}
func (ie *IfExpression) SetType(t types.Type) types.Type {
	ie.Type = t
	return t
}

// fn (<parameters>) { <body> }
type FunctionExpression struct {
	Token      tokens.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
	Type       types.Type
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) TokenLiteral() string {
	return fe.Token.Literal
}
func (fe *FunctionExpression) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fe.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fe.Body.String())

	return out.String()
}
func (fe *FunctionExpression) GetType() types.Type {
	return fe.Type
}
func (fe *FunctionExpression) SetType(t types.Type) types.Type {
	fe.Type = t
	return t
}

// <identifier>(<arguments>)
type CallExpression struct {
	Token      tokens.Token
	Identifier *IdentifierExpression
	Arguments  []Expression
	Type       types.Type
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Identifier.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (ce *CallExpression) GetType() types.Type {
	return ce.Type
}
func (ce *CallExpression) SetType(t types.Type) types.Type {
	ce.Type = t
	return t
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

// while ( <condition> ) { <body> }
type WhileStatement struct {
	Token     tokens.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}

// for ( ?(<initialization>); ?(<condition>); ?(<update>) ) { <body> }
// ?(...) = optional
type ForStatement struct {
	Token          tokens.Token
	Initialization Statement
	Condition      Expression
	Update         Expression
	Body           *BlockStatement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString("(")

	if fs.Initialization != nil {
		out.WriteString(fs.Initialization.String())
	}

	out.WriteString("; ")

	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}

	out.WriteString("; ")

	if fs.Update != nil {
		out.WriteString(fs.Update.String())
	}

	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	out.WriteString("}")

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
