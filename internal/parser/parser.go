package parser

import (
	"fmt"
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/tokens"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	LOWEST = iota
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	UNARY
	FUNCTION_CALL
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
	nextToken    tokens.Token
	errors       []string

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
	}

	p.readToken()
	p.readToken()

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)

	p.registerPrefix(tokens.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(tokens.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(tokens.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(tokens.BANG, p.parseUnaryExpression)
	p.registerPrefix(tokens.MINUS, p.parseUnaryExpression)

	return p
}

// utils
func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) isCurrentToken(tt tokens.TokenType) bool {
	return p.currentToken.Type == tt
}

func (p *Parser) isNextToken(tt tokens.TokenType) bool {
	return p.nextToken.Type == tt
}

func (p *Parser) expectPeek(tt tokens.TokenType) bool {
	if p.isNextToken(tt) {
		p.readToken()
		return true
	} else {
		p.peekErrorBuilder(tt)
		return false
	}
}

func (p *Parser) readTill(tt tokens.TokenType) {
	for !p.isCurrentToken(tt) {
		p.readToken()
	}
}

func (p *Parser) registerPrefix(tt tokens.TokenType, prefixParseFn prefixParseFn) {
	p.prefixParseFns[tt] = prefixParseFn
}

func (p *Parser) registerInfix(tt tokens.TokenType, infixParseFn infixParseFn) {
	p.infixParseFns[tt] = infixParseFn
}

// error builders
func (p *Parser) peekErrorBuilder(tt tokens.TokenType) {
	msg := fmt.Sprintf("[line %d] expected next token to be %s, got %s instead", p.nextToken.Line, tt, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnErrorBuilder(tt tokens.TokenType) {
	msg := fmt.Sprintf("[line %d] no prefix parse fn for %s token found", p.nextToken.Line, tt)
	p.errors = append(p.errors, msg)
}

// expression parsers
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	i, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: i,
	}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	f, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		return nil
	}

	return &ast.FloatLiteral{
		Token: p.currentToken,
		Value: f,
	}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	stmt := &ast.UnaryExpression{
		Token: p.currentToken,
	}

	p.readToken()
	stmt.Expression = p.parseExpression(UNARY)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnErrorBuilder(p.currentToken.Type)
		return nil
	}

	exp := prefix()
	return exp
}

// statement parsers
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(tokens.IDENTIFIER) {
		return nil
	}

	stmt.Identifier = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}

	p.readTill(tokens.SEMICOLON)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.readToken()

	p.readTill(tokens.SEMICOLON)
	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// public methods
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.readToken()
	}

	return program
}
