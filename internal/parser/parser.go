package parser

import (
	"slices"
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
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
	tokens.ASSIGN:              ASSIGN,
	tokens.EQUALS:              EQUALS,
	tokens.NOT_EQUALS:          EQUALS,
	tokens.LESS_THAN:           LESS_GREATER,
	tokens.GREATER_THAN:        LESS_GREATER,
	tokens.LESS_THAN_EQUALS:    LESS_GREATER,
	tokens.GREATER_THAN_EQUALS: LESS_GREATER,
	tokens.PLUS:                SUM,
	tokens.MINUS:               SUM,
	tokens.STAR:                PRODUCT,
	tokens.SLASH:               PRODUCT,
	tokens.LPAREN:              FUNCTION_CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	tokens    []tokens.Token
	nextIdx   int
	currToken tokens.Token
	nextToken tokens.Token
	errors    []string

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(tks []tokens.Token) *Parser {
	p := &Parser{
		tokens:  tks,
		nextIdx: -1,
	}
	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)

	p.registerPrefixFn(tokens.IDENTIFIER, p.parseIdentifierExpression)
	p.registerPrefixFn(tokens.STRING, p.parseStringExpression)
	p.registerPrefixFn(tokens.INTEGER, p.parseIntegerExpression)
	p.registerPrefixFn(tokens.TRUE, p.parseBooleanExpression)
	p.registerPrefixFn(tokens.FALSE, p.parseBooleanExpression)
	p.registerPrefixFn(tokens.FLOAT, p.parseFloatExpression)
	p.registerPrefixFn(tokens.MINUS, p.parseUnaryExpression)
	p.registerPrefixFn(tokens.BANG, p.parseUnaryExpression)
	p.registerPrefixFn(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFn(tokens.IF, p.parseIfExpression)
	p.registerPrefixFn(tokens.FUNCTION, p.parseFunctionExpression)

	p.registerInfixFn(tokens.PLUS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.MINUS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.STAR, p.parseBinaryExpression)
	p.registerInfixFn(tokens.SLASH, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LESS_THAN, p.parseBinaryExpression)
	p.registerInfixFn(tokens.GREATER_THAN, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LESS_THAN_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.GREATER_THAN_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.NOT_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LPAREN, p.parseCallExpression)

	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) registerPrefixFn(tt tokens.TokenType, prefixFn prefixParseFn) {
	p.prefixParseFns[tt] = prefixFn
}

func (p *Parser) registerInfixFn(tt tokens.TokenType, infixFn infixParseFn) {
	p.infixParseFns[tt] = infixFn
}

func (p *Parser) addError(err string) {
	p.errors = append(p.errors, err)
}

func (p *Parser) readToken() {
	p.currToken = p.nextToken
	p.nextIdx++

	if p.nextIdx >= len(p.tokens) {
		p.nextToken = tokens.New(tokens.EOF, "", p.currToken.Line, p.currToken.StartColumn, p.currToken.EndColumn)
	} else {
		p.nextToken = p.tokens[p.nextIdx]
	}
}

func (p *Parser) peekToken() tokens.Token {
	if p.nextIdx >= len(p.tokens) {
		return tokens.New(tokens.EOF, "", p.currToken.Line, p.currToken.StartColumn, p.currToken.EndColumn)
	} else {
		return p.tokens[p.nextIdx]
	}
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedenceTable[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceTable[p.peekToken().Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) isCurrentToken(tt tokens.TokenType) bool {
	return p.currToken.Type == tt
}

func (p *Parser) isNextToken(tt tokens.TokenType) bool {
	nextToken := p.peekToken()
	return nextToken.Type == tt
}

func (p *Parser) checkAndReadToken(tt tokens.TokenType) bool {
	nextToken := p.peekToken()

	if nextToken.Type == tt {
		p.readToken()
		return true
	} else {
		p.addError(utils.ParserExpectedNextTokenToBeErrorBuilder(nextToken, tt))
		return false
	}
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.IdentifierExpression{
		Token:   p.currToken,
		Literal: p.currToken.Literal,
	}
}

func (p *Parser) parseStringExpression() ast.Expression {
	return &ast.StringExpression{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseIntegerExpression() ast.Expression {
	v, err := strconv.ParseInt(p.currToken.Literal, 10, 64)
	if err != nil {
		p.addError(utils.ParserFailedToParseExpressionErrorBuilder(p.currToken, err.Error()))
		return nil
	}

	return &ast.IntegerExpression{
		Token: p.currToken,
		Value: v,
	}
}

func (p *Parser) parseFloatExpression() ast.Expression {
	v, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil {
		p.addError(utils.ParserFailedToParseExpressionErrorBuilder(p.currToken, err.Error()))
		return nil
	}

	return &ast.FloatExpression{
		Token: p.currToken,
		Value: v,
	}
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	isTrue := p.currToken.Type == tokens.TRUE

	return &ast.BooleanExpression{
		Token: p.currToken,
		Value: isTrue,
	}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	unaryOperator := p.currToken
	expr := &ast.UnaryExpression{
		Token: p.currToken,
	}
	p.readToken()

	expr.Expr = p.parseExpression(LOWEST)
	if expr.Expr == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(unaryOperator))
		return nil
	}

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	expr := &ast.GroupedExpression{}

	lParenToken := p.currToken
	// consume left paren token
	p.readToken()
	expr.Expr = p.parseExpression(LOWEST)

	if expr.Expr == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(lParenToken))
		return nil
	}

	if !p.checkAndReadToken(tokens.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.currToken,
	}
	if !p.checkAndReadToken(tokens.LPAREN) {
		return nil
	}

	expr.Condition = p.parseExpression(LOWEST)
	if expr.Condition == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(expr.Token))
		return nil
	}

	if !p.checkAndReadToken(tokens.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.isNextToken(tokens.ELSE) {
		// land on else token
		p.readToken()

		if p.isNextToken(tokens.IF) {
			// land on if token
			p.readToken()

			expr.Alternative = &ast.BlockStatement{
				Token: p.currToken,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: p.currToken,
						Expr:  p.parseIfExpression(),
					},
				},
			}
		} else {
			if !p.checkAndReadToken(tokens.LBRACE) {
				return nil
			}

			expr.Alternative = p.parseBlockStatement()
		}
	}

	return expr
}

func (p *Parser) parseFunctionParameters() []*ast.IdentifierExpression {
	parameters := []*ast.IdentifierExpression{}

	if p.isNextToken(tokens.RPAREN) {
		p.readToken() // consume left paren
		return parameters
	}

	p.readToken()

	parameters = append(parameters, &ast.IdentifierExpression{
		Token:   p.currToken,
		Literal: p.currToken.Literal,
	})

	for p.isNextToken(tokens.COMMA) {
		p.readToken() // consume previous parameter
		p.readToken() // consume comma

		parameters = append(parameters, &ast.IdentifierExpression{
			Token:   p.currToken,
			Literal: p.currToken.Literal,
		})
	}

	if !p.checkAndReadToken(tokens.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	expr := &ast.FunctionExpression{
		Token: p.currToken,
	}

	if !p.checkAndReadToken(tokens.LPAREN) {
		return nil
	}
	expr.Parameters = p.parseFunctionParameters()

	if !p.checkAndReadToken(tokens.LBRACE) {
		return nil
	}
	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.isNextToken(tokens.RPAREN) {
		p.readToken() // consume left paren
		return args
	}

	p.readToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.isNextToken(tokens.COMMA) {
		p.readToken() // consume previous argument
		p.readToken() // consume comma

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.checkAndReadToken(tokens.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token: p.currToken,
	}

	identifier, ok := left.(*ast.IdentifierExpression)
	if !ok {
		p.addError(utils.ParseExpectedXExpressionErrorBuilder[*ast.IdentifierExpression](p.currToken, left))
		return nil
	}

	expr.Identifier = identifier
	expr.Arguments = p.parseCallArguments()

	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Left:     left,
		Operator: p.currToken,
	}

	precedence := p.currentPrecedence()
	binaryOperator := p.currToken
	p.readToken()

	expr.Right = p.parseExpression(precedence)
	if expr.Right == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(binaryOperator))
		return nil
	}

	return expr
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.addError(utils.ParserNoPrefixFnErrorBuilder(p.currToken))
		return nil
	}

	expr := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken().Type]
		if infix == nil {
			return expr
		}

		p.readToken()
		expr = infix(expr)
	}

	return expr
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

	if p.isNextToken(tokens.ASSIGN) {
		p.readToken()
		assignToken := p.currToken
		p.readToken()

		stmt.Value = p.parseExpression(LOWEST)
		if stmt.Value == nil {
			p.addError(utils.ParserExpressionExpectedErrorBuilder(assignToken))
			return nil
		}
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
		p.addError(utils.ParserExpressionExpectedErrorBuilder(returnToken))
		return nil
	}

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
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
		p.addError(utils.ParserExpressionExpectedErrorBuilder(lParenToken))
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.WHILE:
		return p.parseWhileStatement()
	case tokens.FOR:
		return p.parseForStatement()
	case tokens.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.readToken()
	}

	return program
}

func (p *Parser) Error() []string {
	slices.Reverse(p.errors)
	return p.errors
}
