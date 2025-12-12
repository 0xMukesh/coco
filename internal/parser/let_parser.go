package parser

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

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
