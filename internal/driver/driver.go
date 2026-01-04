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
	tokens []tokens.Token
	ast    ast.Node
}

func NewDriver(source string) *Driver {
	return &Driver{
		Source: source,
	}
}

func (d *Driver) Lex() ([]tokens.Token, error) {
	l := lexer.New(d.Source)
	d.tokens = []tokens.Token{}

	for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
		d.tokens = append(d.tokens, tok)
	}

	// FIXME: handle lexing errors over here
	return d.tokens, nil
}

func (d *Driver) Parse() (ast.Node, error) {
	if d.tokens == nil {
		return nil, errors.New("source not lexed: tokens are nil")
	}

	p := parser.New(d.tokens)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("failed to parse program -  %s", strings.Join(p.Errors(), "\n"))
	}

	d.ast = program
	return program, nil
}

func (d *Driver) Eval() (object.Object, error) {
	if d.ast == nil {
		return nil, errors.New("ast not found")
	}

	res := eval.Eval(d.ast)
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
