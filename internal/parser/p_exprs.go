package parser

import (
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.addError(ParserNoPrefixFnErrorBuilder(p.currToken))
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
		p.addError(ParserFailedToParseExpressionErrorBuilder(p.currToken, err.Error()))
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
		p.addError(ParserFailedToParseExpressionErrorBuilder(p.currToken, err.Error()))
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

	expr.Expr = p.parseExpression(UNARY)
	if expr.Expr == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(unaryOperator))
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
		p.addError(ParserExpressionExpectedErrorBuilder(lParenToken))
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
		p.addError(ParserExpressionExpectedErrorBuilder(expr.Token))
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
		p.addError(ParseExpectedXExpressionErrorBuilder[*ast.IdentifierExpression](p.currToken, left))
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

	// exponentiation is right associative
	if binaryOperator.Type == tokens.DOUBLE_STAR {
		expr.Right = p.parseExpression(precedence - 1)
	} else {
		expr.Right = p.parseExpression(precedence)
	}

	if expr.Right == nil {
		p.addError(ParserExpressionExpectedErrorBuilder(binaryOperator))
		return nil
	}

	return expr
}
