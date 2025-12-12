package ast

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
