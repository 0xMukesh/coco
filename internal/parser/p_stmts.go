package parser

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.BREAK:
		return p.parseBreakStatement()
	case tokens.WHILE:
		return p.parseWhileStatement()
	case tokens.FOR:
		return p.parseForStatement()
	case tokens.LBRACE:
		return p.parseBlockStatement()
	case tokens.IDENTIFIER:
		if p.isNextToken(tokens.ASSIGN) {
			return p.parseAssignmentStatement()
		}

		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.checkAndReadToken(tokens.IDENTIFIER) {
		return nil
	}

	stmt.Identifier = &ast.IdentifierExpression{
		Token:   p.currToken,
		Literal: p.currToken.Literal,
	}

	if !p.checkAndReadToken(tokens.ASSIGN) {
		return nil
	}

	assignToken := p.currToken
	p.readToken()

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(assignToken))
		return nil
	}

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{
		Token: p.currToken,
		Identifier: &ast.IdentifierExpression{
			Token:   p.currToken,
			Literal: p.currToken.Literal,
		},
	}

	if !p.checkAndReadToken(tokens.ASSIGN) {
		return nil
	}

	assignToken := p.currToken
	p.readToken()

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(assignToken))
		return nil
	}

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.currToken,
	}
	returnToken := p.currToken
	p.readToken()

	stmt.Expr = p.parseExpression(LOWEST)
	if stmt.Expr == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(returnToken))
		return nil
	}

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{
		Token: p.currToken,
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.currToken,
	}

	if !p.checkAndReadToken(tokens.LPAREN) {
		return nil
	}
	lParenToken := p.currToken

	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(lParenToken))
		return nil
	}

	if !p.checkAndReadToken(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{
		Token: p.currToken,
	}
	if !p.checkAndReadToken(tokens.LPAREN) {
		return nil
	}
	// consume left paren
	p.readToken()

	if !p.isCurrentToken(tokens.SEMICOLON) {
		// if initialization statement is not empty, then parse it
		stmt.Initialization = p.parseStatement()
	}

	// consume the semicolon after initialization statement
	p.readToken()

	if !p.isCurrentToken(tokens.SEMICOLON) {
		// if condition expression is not empty, then parse it
		stmt.Condition = p.parseExpression(LOWEST)

		// check if there is a semicolon after the condition expression
		if !p.checkAndReadToken(tokens.SEMICOLON) {
			return nil
		}
	}

	// consume the semicolon after condition expression
	p.readToken()

	if !p.isCurrentToken(tokens.RPAREN) {
		// if update expression is not empty, then parse it
		stmt.Update = p.parseExpression(LOWEST)
		// consume right paren token
		p.readToken()
	}

	if !p.checkAndReadToken(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currToken,
	}
	block.Statements = []ast.Statement{}

	p.readToken() // consume LBRACE
	for !p.isCurrentToken(tokens.RBRACE) && !p.isCurrentToken(tokens.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.readToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}
	stmt.Expr = p.parseExpression(LOWEST)

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}
