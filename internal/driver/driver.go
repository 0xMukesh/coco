package driver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/eval"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/object"
	"github.com/0xmukesh/coco/internal/parser"
	"github.com/0xmukesh/coco/internal/tokens"
)

type Driver struct {
	Source string
	Tokens []tokens.Token
	Ast    ast.Node
}

func NewDriver(source string) *Driver {
	return &Driver{
		Source: source,
	}
}

func (d *Driver) Lex() ([]tokens.Token, error) {
	l := lexer.New(d.Source)
	d.Tokens = []tokens.Token{}

	for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
		d.Tokens = append(d.Tokens, tok)
	}

	// FIXME: handle lexing errors over here
	return d.Tokens, nil
}

func (d *Driver) Parse() (ast.Node, error) {
	if d.Tokens == nil {
		return nil, errors.New("source not lexed: tokens are nil")
	}

	p := parser.New(d.Tokens)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("failed to parse program -  %s", strings.Join(p.Errors(), "\n"))
	}

	d.Ast = program
	return program, nil
}

func (d *Driver) Eval() (object.Object, error) {
	if d.Ast == nil {
		return nil, errors.New("ast not found")
	}

	env := object.NewEnvironment()
	evaler := eval.NewEvalutor()

	res := evaler.Eval(d.Ast, env)
	return res, nil
}

func (d *Driver) Process() (object.Object, error) {
	_, err := d.Lex()
	if err != nil {
		return nil, err
	}

	_, err = d.Parse()
	if err != nil {
		return nil, err
	}

	res, err := d.Eval()
	if err != nil {
		return nil, err
	}

	return res, nil
}
