package parser

import (
	"fmt"
	"strconv"

	"github.com/udeshyadhungana/interprerer/app/ast"
	"github.com/udeshyadhungana/interprerer/app/lexer"
	"github.com/udeshyadhungana/interprerer/app/token"
)

/* for pratt's */
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	/* For pratt's parser */
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

const (
	_ int = iota
	LOWEST
	ASSIGN
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: CALL,
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.JHUT_MUJI, p.parseBoolean)
	p.registerPrefix(token.SACHO_MUJI, p.parseBoolean)
	p.registerPrefix(token.KAAM_GAR_MUJI, p.parseKaamGarMuji)
	p.registerPrefix(token.LPAREN, p.parseLeftParenthesis)
	p.registerPrefix(token.YEDI_MUJI, p.parseYediMujiExpression)
	p.registerPrefix(token.STRING, p.parseStringExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayExpression)
	p.registerPrefix(token.JABA_SAMMA_MUJI, p.parseJabasammaMujiExpression)
	// p.registerPrefix(token.ASSIGN, p.parseReassignmentExpression)

	// infix functions for operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseArrayIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next tokenn to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{Token: p.curToken}
	s := []ast.Statement{}
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}
	p.nextToken()
	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		s = append(s, stmt)
		p.nextToken()
	}
	bs.Statements = s
	return bs
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.THOOS_MUJI:
		return p.parseThoosMujiStatement()
	case token.PATHA_MUJI:
		return p.parsePathaMujiStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseThoosMujiStatement() *ast.ThoosMujiStatement {
	stmt := &ast.ThoosMujiStatement{Token: p.curToken}

	if !p.expectPeek(token.IDFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if !p.curTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, "expected semicolon at the end of statement")
		return nil
	}

	return stmt
}

func (p *Parser) parsePathaMujiStatement() *ast.PathaMujiStatement {
	stmt := &ast.PathaMujiStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if !p.curTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, "expected semicolon at the end of statement")
	}

	return stmt
}

func (p *Parser) parseYediMujiExpression() ast.Expression {
	stmt := &ast.YediMujiExpression{Token: p.curToken}
	if !p.curTokenIs(token.YEDI_MUJI) {
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	// evaluate expression
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	if !p.curTokenIs(token.RPAREN) && !p.peekTokenIs(token.LBRACE) {
		return nil
	}
	p.nextToken()

	stmt.Consequent = p.parseBlockStatement()

	if !p.peekTokenIs(token.NABHAE_CHIKNE) {
		return stmt
	}
	p.nextToken()
	p.nextToken()
	stmt.Alternative = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseJabasammaMujiExpression() ast.Expression {
	stmt := &ast.JabasammaMujiExpression{Token: p.curToken}
	if !p.curTokenIs(token.JABA_SAMMA_MUJI) {
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	// parse condition
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	if !p.curTokenIs(token.RPAREN) && !p.peekTokenIs(token.LBRACE) {
		return nil
	}
	p.nextToken()
	stmt.Consequent = p.parseBlockStatement()
	return stmt
}

/* expressions */
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// the most beautiful function of this program
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	switch p.curToken.Type {
	case token.SACHO_MUJI:
		return &ast.Boolean{Token: p.curToken, Value: true}
	case token.JHUT_MUJI:
		return &ast.Boolean{Token: p.curToken, Value: false}
	default:
		panic(fmt.Sprintf("not a boolean: %s\n", p.curToken.Literal))
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for (%s) found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseLeftParenthesis() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	result := ast.CallExpression{Token: p.curToken, Function: function}
	result.Arguments = p.parseArguments()
	return &result
}

func (p *Parser) parseArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseKaamGarMuji() ast.Expression {
	result := ast.KaamGarMujiExpression{Token: p.curToken}
	if !p.curTokenIs(token.KAAM_GAR_MUJI) {
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	p.nextToken()
	if p.curTokenIs(token.RPAREN) {
		result.Arguments = nil
	} else {
		for {
			currentArg := p.parseIdentifier()
			r, ok := currentArg.(*ast.Identifier)
			if !ok {
				p.errors = append(p.errors, "parameters must be identifiers")
			}
			result.Arguments = append(result.Arguments, r)
			p.nextToken()
			if p.curTokenIs(token.RPAREN) {
				break
			}
			if !p.curTokenIs(token.COMMA) {
				p.errors = append(p.errors, "expected comma after argument")
				return nil
			}
			p.nextToken()
		}
	}
	p.nextToken()
	result.Body = p.parseBlockStatement()
	return &result
}

func (p *Parser) parseStringExpression() ast.Expression {
	return &ast.StringExpression{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayExpression() ast.Expression {
	result := &ast.ArrayExpression{Token: p.curToken}
	var current ast.Expression

	if p.peekTokenIs(token.RBRACKET) {
		return result
	}
	p.nextToken()
	for {
		current = p.parseExpression(LOWEST)
		result.Elements = append(result.Elements, current)
		result.Length += 1
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		} else if p.curTokenIs(token.RBRACKET) {
			break
		} else {
			p.errors = append(p.errors, "malformed array expression")
			return nil
		}
	}
	return result
}

func (p *Parser) parseArrayIndexExpression(arrExpr ast.Expression) ast.Expression {
	result := &ast.ArrayIndexExpression{Token: p.curToken, Array: arrExpr}
	p.nextToken()
	indexExpr := p.parseExpression(LOWEST)
	p.nextToken()
	result.Index = indexExpr
	return result
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
