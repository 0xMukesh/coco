package parser

import (
	"slices"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

const (
	LOWEST         = iota
	ASSIGN         // =
	LOGICAL        // &&, ||
	COMPARISON     // >, >=, <, <=, ==, !=
	ADDITION       // +, -
	MULTIPLICATION // *, /, %
	EXPONENTIATION // **
	UNARY
	FUNCTION_CALL
)

var precedenceTable = map[tokens.TokenType]int{
	tokens.ASSIGN:              ASSIGN,
	tokens.AND:                 LOGICAL,
	tokens.OR:                  LOGICAL,
	tokens.EQUALS:              COMPARISON,
	tokens.NOT_EQUALS:          COMPARISON,
	tokens.LESS_THAN:           COMPARISON,
	tokens.GREATER_THAN:        COMPARISON,
	tokens.LESS_THAN_EQUALS:    COMPARISON,
	tokens.GREATER_THAN_EQUALS: COMPARISON,
	tokens.MINUS:               ADDITION,
	tokens.MINUS_EQUAL:         ADDITION,
	tokens.PLUS:                ADDITION,
	tokens.PLUS_EQUAL:          ADDITION,
	tokens.STAR:                MULTIPLICATION,
	tokens.STAR_EQUAL:          MULTIPLICATION,
	tokens.SLASH:               MULTIPLICATION,
	tokens.SLASH_EQUAL:         MULTIPLICATION,
	tokens.MODULO:              MULTIPLICATION,
	tokens.DOUBLE_STAR:         EXPONENTIATION,
	tokens.INCREMENT:           UNARY,
	tokens.DECREMENT:           UNARY,
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
	p.registerPrefixFn(tokens.INCREMENT, p.parseUnaryExpression)
	p.registerPrefixFn(tokens.DECREMENT, p.parseUnaryExpression)
	p.registerPrefixFn(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFn(tokens.IF, p.parseIfExpression)
	p.registerPrefixFn(tokens.FUNCTION, p.parseFunctionExpression)

	p.registerInfixFn(tokens.PLUS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.MINUS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.STAR, p.parseBinaryExpression)
	p.registerInfixFn(tokens.SLASH, p.parseBinaryExpression)
	p.registerInfixFn(tokens.MODULO, p.parseBinaryExpression)
	p.registerInfixFn(tokens.PLUS_EQUAL, p.parseBinaryExpression)
	p.registerInfixFn(tokens.MINUS_EQUAL, p.parseBinaryExpression)
	p.registerInfixFn(tokens.STAR_EQUAL, p.parseBinaryExpression)
	p.registerInfixFn(tokens.SLASH_EQUAL, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LESS_THAN, p.parseBinaryExpression)
	p.registerInfixFn(tokens.GREATER_THAN, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LESS_THAN_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.GREATER_THAN_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.NOT_EQUALS, p.parseBinaryExpression)
	p.registerInfixFn(tokens.OR, p.parseBinaryExpression)
	p.registerInfixFn(tokens.AND, p.parseBinaryExpression)
	p.registerInfixFn(tokens.DOUBLE_STAR, p.parseBinaryExpression)
	p.registerInfixFn(tokens.LPAREN, p.parseCallExpression)

	p.readToken()
	p.readToken()

	return p
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

func (p *Parser) Errors() []string {
	slices.Reverse(p.errors)
	return p.errors
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
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
		p.addError(ParserExpectedNextTokenToBeErrorBuilder(nextToken, tt))
		return false
	}
}
