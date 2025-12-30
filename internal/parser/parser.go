package parser

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/tokens"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
	nextToken    tokens.Token
	errors       []string
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
	}

	p.readToken()
	p.readToken()

	return p
}

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

func (p *Parser) peekErrorBuilder(tt tokens.TokenType) {
	msg := fmt.Sprintf("[line %d] expected next token to be %s, got %s instead", p.nextToken.Line, tt, p.nextToken.Type)
	p.errors = append(p.errors, msg)
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
		return nil
	}
}

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
