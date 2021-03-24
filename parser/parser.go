package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression
)

const (
	_                  int = iota
	LOWEST                 // 1
	EQUALS                 // ==
	LESSGREATER            // < or >
	SUM                    // +
	PRODUCT                // *
	PREFIX                 // -X or !X
	ARRAY_ACCESS           // array[x]
	CALL                   // function(x)
	EXTERNAL_REFERENCE     //x.y
)

var precedences = map[string]int{
	token.EQUAL:        EQUALS,
	token.NOT_EQUAL:    EQUALS,
	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.DIVIDE:       PRODUCT,
	token.MULTIPLY:     PRODUCT,
	token.LSQBRACKET:   ARRAY_ACCESS,
	token.LPAREN:       CALL,
	token.DOT:          EXTERNAL_REFERENCE,
}

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token

	prefixParserFns map[string]prefixParserFn
	infixParserFns  map[string]infixParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	p.prefixParserFns = make(map[string]prefixParserFn)
	p.infixParserFns = make(map[string]infixParserFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LSQBRACKET, p.parseArray)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.DOT, p.parseExternalReference)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LSQBRACKET, p.parseArrayAccessExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) AddError(expectedType string) {
	newError := fmt.Sprintf("Expected token of type %s, but got %s instead", expectedType, p.peekToken.Type)
	p.errors = append(p.errors, newError)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectToken(tokenType string) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	}
	p.AddError(tokenType)
	return false
}

func (p *Parser) registerPrefix(tokenType string, fn prefixParserFn) {
	p.prefixParserFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType string, fn infixParserFn) {
	p.infixParserFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if number, ok := precedences[p.peekToken.Type]; ok {
		return number
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if number, ok := precedences[p.currentToken.Type]; ok {
		return number
	}
	return LOWEST
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.USE:
		return p.parseUseStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENTIFIER:
		if p.peekToken.Type == token.ASSIGN {
			return p.parseReassignmentStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseUseStatement() ast.Statement {
	statement := &ast.UseStatement{Token: p.currentToken}

	if !p.expectToken(token.IDENTIFIER) {
		return nil
	}

	statement.Filename = p.currentToken.Literal

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: p.currentToken}

	if !p.expectToken(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectToken(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

//TODO: finish this (add support to while loops)
func (p *Parser) parseWhileStatement() ast.Statement {
	statement := &ast.WhileStatement{Token: p.currentToken}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	if p.peekToken.Type == token.RPAREN {
		p.AddError("while loop must have a condition")
		return nil
	}

	p.nextToken()

	statement.Condition = p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	if !p.expectToken(token.LBRACE) {
		return nil
	}

	statement.Block = p.parseBlockStatement()

	return statement
}

func (p *Parser) parseReassignmentStatement() ast.Statement {
	statement := &ast.ReassignmentStatement{Token: p.currentToken}

	statement.Variable = ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	p.nextToken()
	p.nextToken()

	statement.NewValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParserFns[p.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {

		infix := p.infixParserFns[p.peekToken.Type]
		if infix == nil {
			return nil
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseGroupedExpression() ast.Expression {

	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	il := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s to an integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	il.Value = value

	return il
}

func (p *Parser) parseStringLiteral() ast.Expression {
	sl := &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}

	return sl
}

func (p *Parser) parseArray() ast.Expression {
	arr := &ast.Array{Token: p.currentToken}
	arr.Elements = []ast.Expression{}

	if p.peekToken.Type == token.RSQBRACKET {
		p.nextToken()
		return arr
	}

	p.nextToken()

	arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))

	p.nextToken()

	for p.currentToken.Type != token.RSQBRACKET {
		if p.currentToken.Type != token.COMMA {
			p.AddError(token.COMMA)
			return nil
		}
		p.nextToken()
		arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))
		p.nextToken()
	}

	return arr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.currentToken, Left: left, Operator: p.currentToken.Literal}

	precedence := p.currentPrecedence()

	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	expression := &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	if !p.expectToken(token.LBRACE) {
		return nil
	}

	expression.TrueBlock = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if !p.expectToken(token.LBRACE) {
			return nil
		}

		expression.FalseBlock = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() ast.BlockStatement {
	bs := ast.BlockStatement{Token: p.currentToken}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	for p.currentToken.Type != token.RBRACE {
		statement := p.parseStatement()
		if statement != nil {
			bs.Statements = append(bs.Statements, statement)
		}
		p.nextToken()
	}

	return bs
}

func (p *Parser) parseFunctionLiteral() ast.Expression {

	fl := &ast.FunctionLiteral{Token: p.currentToken}
	fl.Parameters = []ast.Identifier{}

	if p.peekToken.Type == token.IDENTIFIER {
		fl.Name = ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal}
		p.nextToken()
	}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	p.nextToken()

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectToken(token.LBRACE) {
		return nil
	}

	fl.Block = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []ast.Identifier {
	identifiers := []ast.Identifier{}

	if p.currentToken.Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	identifiers = append(identifiers, ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal})

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal})
	}

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Function: left}

	expression.Arguments = p.parseCallArguments()

	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {

	arguments := []ast.Expression{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return arguments
	}

	p.nextToken()

	arguments = append(arguments, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	return arguments
}

func (p *Parser) parseArrayAccessExpression(left ast.Expression) ast.Expression {
	arrAccess := &ast.ArrayAccessExpression{Token: p.currentToken, Array: left}

	p.nextToken()

	arrAccess.Position = p.parseExpression(LOWEST)

	if !p.expectToken(token.RSQBRACKET) {
		return nil
	}

	return arrAccess
}

func (p *Parser) parseExternalReference(left ast.Expression) ast.Expression {
	expression := &ast.ExternalReferenceExpression{Token: p.currentToken}

	expression.Module = left.String()

	if p.peekToken.Type != token.IDENTIFIER {
		p.AddError(token.IDENTIFIER)
		return nil
	}

	p.nextToken()

	expression.Referece = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return expression
}
