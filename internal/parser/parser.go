package parser

import (
	"fmt"
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	LOWEST = iota
	ASSIGN
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	UNARY
	FUNCTION_CALL
)

var precedenceTable = map[tokens.TokenType]int{
	tokens.ASSIGN:             ASSIGN,
	tokens.EQUALS:             EQUALS,
	tokens.NOT_EQUALS:         EQUALS,
	tokens.LESS_THAN:          LESS_GREATER,
	tokens.GREATER_THAN:       LESS_GREATER,
	tokens.LESS_THAN_EQUAL:    LESS_GREATER,
	tokens.GREATER_THAN_EQUAL: LESS_GREATER,
	tokens.PLUS:               SUM,
	tokens.MINUS:              SUM,
	tokens.STAR:               PRODUCT,
	tokens.SLASH:              PRODUCT,
	tokens.LPAREN:             FUNCTION_CALL,
}

type Parser struct {
	tokens          []tokens.Token
	currentTokenIdx int
	currentToken    tokens.Token
	nextToken       tokens.Token
	errors          []string

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(tks []tokens.Token) *Parser {
	p := &Parser{
		tokens:          tks,
		currentTokenIdx: -1,
		currentToken:    tokens.Token{Type: tokens.EOF},
		nextToken:       tokens.Token{Type: tokens.EOF},
	}

	p.readToken()
	p.readToken()

	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)

	p.registerPrefix(tokens.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(tokens.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(tokens.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.BANG, p.parseUnaryExpression)
	p.registerPrefix(tokens.MINUS, p.parseUnaryExpression)
	p.registerPrefix(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(tokens.IF, p.parseIfExpression)
	p.registerPrefix(tokens.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(tokens.EQUALS, p.parseBinaryExpression)
	p.registerInfix(tokens.NOT_EQUALS, p.parseBinaryExpression)
	p.registerInfix(tokens.LESS_THAN, p.parseBinaryExpression)
	p.registerInfix(tokens.GREATER_THAN, p.parseBinaryExpression)
	p.registerInfix(tokens.LESS_THAN_EQUAL, p.parseBinaryExpression)
	p.registerInfix(tokens.GREATER_THAN_EQUAL, p.parseBinaryExpression)
	p.registerInfix(tokens.PLUS, p.parseBinaryExpression)
	p.registerInfix(tokens.MINUS, p.parseBinaryExpression)
	p.registerInfix(tokens.STAR, p.parseBinaryExpression)
	p.registerInfix(tokens.SLASH, p.parseBinaryExpression)
	p.registerInfix(tokens.LPAREN, p.parseCallExpression)
	p.registerInfix(tokens.ASSIGN, p.parseAssignmentExpression)

	return p
}

// utils
func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.currentTokenIdx++

	if p.currentTokenIdx < len(p.tokens) {
		p.nextToken = p.tokens[p.currentTokenIdx]
	} else {
		p.nextToken = tokens.Token{Type: tokens.EOF}
	}
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceTable[p.nextToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedenceTable[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.currentToken,
		Value: p.isCurrentToken(tokens.TRUE),
	}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expr := &ast.UnaryExpression{
		Token: p.currentToken,
	}

	p.readToken()
	expr.Expression = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.readToken() // consume LPAREN token
	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token: p.currentToken,
		Left:  left,
	}

	precedence := p.currentPrecedence()
	p.readToken() // consume the operator
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.currentToken,
	}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	expr.Condition = p.parseExpression(LOWEST) // RPAREN missing case is handled by parseGroupExpression
	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.isNextToken(tokens.ELSE) {
		p.readToken()
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.isNextToken(tokens.RPAREN) {
		p.readToken()
		return identifiers
	}

	p.readToken()
	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.isNextToken(tokens.COMMA) {
		p.readToken()
		if !p.expectPeek(tokens.IDENTIFIER) {
			return nil
		}

		ident := &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifiers = append(identifiers, ident)

	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expr := &ast.FunctionLiteral{
		Token: p.currentToken,
	}
	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	expr.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) isValidCallArgument() bool {
	return p.isNextToken(tokens.IDENTIFIER) || p.isNextToken(tokens.INTEGER) || p.isNextToken(tokens.FLOAT) || p.isNextToken(tokens.FUNCTION)
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.isNextToken(tokens.RPAREN) {
		p.readToken()
		return args
	}

	if !p.isValidCallArgument() {
		return nil
	}

	p.readToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.isNextToken(tokens.COMMA) {
		p.readToken()
		if !p.isValidCallArgument() {
			return nil
		}
		p.readToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	return args
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}
	expr.Arguments = p.parseCallArguments()

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseAssignmentExpression(lhs ast.Expression) ast.Expression {
	expr := &ast.AssignmentExpression{
		Token:      p.currentToken,
		Identifier: lhs,
	}
	p.readToken()

	rhs := p.parseExpression(LOWEST)
	expr.Value = rhs
	return expr
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

	expr := prefix()

	for !p.isNextToken(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return expr
		}

		p.readToken()
		expr = infix(expr)
	}

	return expr
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
	p.readToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.readToken()

	stmt.Expression = p.parseExpression(LOWEST)

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.LBRACE:
		return p.parseBlockStatement()
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
