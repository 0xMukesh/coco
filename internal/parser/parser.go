package parser

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

type Parser struct {
	tokens          []tokens.Token
	currentPosition int
	peekPosition    int
	currTok         *tokens.Token
	errors          []string
}

func New(tokens []tokens.Token) *Parser {
	p := &Parser{tokens: tokens}
	p.readToken()

	return p
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}

	for p.currTok.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.readToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) readToken() *tokens.Token {
	if p.peekPosition >= len(p.tokens) {
		p.currTok = nil
	} else {
		p.currTok = &p.tokens[p.peekPosition]
	}

	p.currentPosition = p.peekPosition
	p.peekPosition++
	return p.currTok
}

func (p *Parser) peekToken() *tokens.Token {
	if p.peekPosition >= len(p.tokens) {
		return nil
	} else {
		return &p.tokens[p.peekPosition]
	}
}

func (p *Parser) currTokenIs(t tokens.TokenType) bool {
	return p.tokens[p.currentPosition].Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	if p.peekPosition >= len(p.tokens) {
		return false
	}

	return p.tokens[p.peekPosition].Type == t
}

func (p *Parser) expectPeek(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.readToken()
		return true
	} else {
		peekTok := p.peekToken()
		p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead", t, peekTok.Literal))
		return false
	}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currTok.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	if p.currTok == nil {
		return nil
	}

	stmt := &ast.LetStatement{Token: *p.currTok}
	if !p.expectPeek(tokens.IDENTIFIER) {
		return nil
	}

	stmt.Identifier = &ast.IdentifierExpression{Token: *p.currTok, Value: p.currTok.Literal}
	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}

	for !p.currTokenIs(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if p.currTok == nil {
		return nil
	}

	stmt := &ast.ReturnStatement{Token: *p.currTok}
	if !p.expectPeek(tokens.IDENTIFIER) {
		return nil
	}

	stmt.Value = &ast.IdentifierExpression{Token: *p.currTok, Value: p.currTok.Literal}

	return stmt
}
