package parser

import (
	"slices"
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type (
	prefixParseFn func() ast.Expression
)

type Parser struct {
	tokens    []tokens.Token
	nextIdx   int
	currToken tokens.Token
	nextToken tokens.Token
	errors    []string

	prefixParseFns map[tokens.TokenType]prefixParseFn
}

func New(tks []tokens.Token) *Parser {
	p := &Parser{
		tokens:  tks,
		nextIdx: -1,
	}
	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)

	p.registerPrefixFn(tokens.IDENTIFIER, p.parseIdentifierExpression)
	p.registerPrefixFn(tokens.INTEGER, p.parseIntegerExpression)
	p.registerPrefixFn(tokens.FLOAT, p.parseFloatExpression)
	p.registerPrefixFn(tokens.MINUS, p.parseUnaryExpression)
	p.registerPrefixFn(tokens.BANG, p.parseUnaryExpression)

	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) registerPrefixFn(tt tokens.TokenType, prefixFn prefixParseFn) {
	p.prefixParseFns[tt] = prefixFn
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
		p.addError(utils.ParserPeekCheckFailErrorBuilder(nextToken, tt))
		return false
	}
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.IdentifierExpression{
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

func (p *Parser) parseUnaryExpression() ast.Expression {
	unaryOperator := p.currToken
	expr := &ast.UnaryExpression{
		Token: p.currToken,
	}
	p.readToken()

	expr.Expr = p.parseExpression()
	if expr.Expr == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(unaryOperator))
		return nil
	}

	return expr
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.addError(utils.ParserNoPrefixFnErrorBuilder(p.currToken))
		return nil
	}

	expr := prefix()

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
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if p.isNextToken(tokens.ASSIGN) {
		p.readToken()
		assignToken := p.currToken
		p.readToken()

		stmt.Value = p.parseExpression()
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

	stmt.Expr = p.parseExpression()
	if stmt.Expr == nil {
		p.addError(utils.ParserExpressionExpectedErrorBuilder(returnToken))
		return nil
	}

	if p.isNextToken(tokens.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currToken,
	}
	stmt.Expr = p.parseExpression()

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
