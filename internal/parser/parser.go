// Package parser реализует Pratt parser (operator-precedence parser) -
// построение AST из потока токенов с корректным учётом приоритета
// и ассоциативности операторов без необходимости явно описывать грамматику
// в виде BNF-правил для каждого уровня приоритета.
package parser

import (
	"fmt"
	"strconv"

	"github.com/emityoffwhite/go-calculator/internal/ast"
	"github.com/emityoffwhite/go-calculator/internal/lexer"
)

// Приоритеты операторов: чем выше число, тем сильнее оператор связывает операнды.
// ^ имеет наивысший приоритет среди бинарных и право-ассоциативен (2^3^2 = 2^(3^2)),
// что обрабатывается отдельно в parseExpression.
const (
	_ int = iota
	LOWEST
	ASSIGNMENT // =
	SUM        // + -
	PRODUCT    // * / %
	PREFIX     // -x +x (унарные)
	POWER      // ^
	CALL       // myFunction(x)
)

var precedences = map[lexer.TokenType]int{
	lexer.ASSIGN:   ASSIGNMENT,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.PERCENT:  PRODUCT,
	lexer.CARET:    POWER,
	lexer.LPAREN:   CALL,
}

// Parser строит AST из токенов, выдаваемых lexer.Lexer.
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token
}

// New создаёт Parser для переданной входной строки.
func New(input string) (*Parser, error) {
	p := &Parser{l: lexer.New(input)}

	// Заполняем curToken и peekToken двумя последовательными чтениями -
	// классический паттерн для парсеров с lookahead в один токен.
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	tok, err := p.l.NextToken()
	if err != nil {
		return err
	}
	p.peekToken = tok
	return nil
}

// ParseExpression разбирает входную строку целиком и возвращает корневой узел AST.
// Возвращает ошибку, если после выражения остались необработанные токены
// (например "2 + 3 4" - корректное "2 + 3", но дальше идёт мусор).
func (p *Parser) ParseExpression() (ast.Node, error) {
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.curToken.Type != lexer.EOF {
		return nil, fmt.Errorf("unexpected token %q at position %d", p.curToken.Literal, p.curToken.Pos)
	}

	return expr, nil
}

// parseExpression - сердце Pratt-парсера. Сначала разбирает один операнд
// (число, переменную, унарный минус, скобки или вызов функции), затем в цикле
// "поглощает" следующие бинарные операторы, чей приоритет выше minPrecedence,
// рекурсивно разбирая правую часть. Это и даёт корректную работу приоритетов
// без отдельной грамматической функции на каждый уровень (parseSum, parseProduct, ...).
func (p *Parser) parseExpression(minPrecedence int) (ast.Node, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for {
		prec, ok := precedences[p.curToken.Type]
		if !ok || prec <= minPrecedence {
			break
		}

		op := p.curToken
		if err := p.nextToken(); err != nil {
			return nil, err
		}

		// ^ право-ассоциативен: разбираем правую часть с приоритетом POWER-1,
		// чтобы 2^3^2 стало 2^(3^2), а не (2^3)^2.
		nextMin := prec
		if op.Type == lexer.CARET {
			nextMin = prec - 1
		}

		right, err := p.parseExpression(nextMin)
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpr{Left: left, Op: op.Literal, Right: right}
	}

	return left, nil
}

// parsePrefix разбирает один "атом" выражения: число, идентификатор
// (переменную, вызов функции или присваивание), унарный минус/плюс, либо
// выражение в скобках.
func (p *Parser) parsePrefix() (ast.Node, error) {
	switch p.curToken.Type {
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.IDENT:
		return p.parseIdentifierOrCall()
	case lexer.MINUS, lexer.PLUS:
		return p.parseUnary()
	case lexer.LPAREN:
		return p.parseGrouped()
	default:
		return nil, fmt.Errorf("unexpected token %q at position %d", p.curToken.Literal, p.curToken.Pos)
	}
}

func (p *Parser) parseNumber() (ast.Node, error) {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q at position %d", p.curToken.Literal, p.curToken.Pos)
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	return &ast.NumberLiteral{Value: value}, nil
}

func (p *Parser) parseUnary() (ast.Node, error) {
	op := p.curToken.Literal
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	operand, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return &ast.UnaryExpr{Op: op, Operand: operand}, nil
}

func (p *Parser) parseGrouped() (ast.Node, error) {
	if err := p.nextToken(); err != nil { // пропускаем '('
		return nil, err
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.curToken.Type != lexer.RPAREN {
		return nil, fmt.Errorf("expected ')' at position %d, got %q", p.curToken.Pos, p.curToken.Literal)
	}
	if err := p.nextToken(); err != nil { // пропускаем ')'
		return nil, err
	}

	return expr, nil
}

// parseIdentifierOrCall разбирает три возможных случая для идентификатора:
//  1. имя(...)  - вызов функции
//  2. имя = выражение - присваивание переменной
//  3. имя - просто чтение переменной
func (p *Parser) parseIdentifierOrCall() (ast.Node, error) {
	name := p.curToken.Literal
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	switch p.curToken.Type {
	case lexer.LPAREN:
		return p.parseCall(name)
	case lexer.ASSIGN:
		if err := p.nextToken(); err != nil { // пропускаем '='
			return nil, err
		}
		value, err := p.parseExpression(ASSIGNMENT)
		if err != nil {
			return nil, err
		}
		return &ast.AssignExpr{Name: name, Value: value}, nil
	default:
		return &ast.Identifier{Name: name}, nil
	}
}

func (p *Parser) parseCall(function string) (ast.Node, error) {
	if err := p.nextToken(); err != nil { // пропускаем '('
		return nil, err
	}

	var args []ast.Node

	if p.curToken.Type == lexer.RPAREN {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.CallExpr{Function: function, Args: args}, nil
	}

	for {
		arg, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.curToken.Type == lexer.COMMA {
			if err := p.nextToken(); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	if p.curToken.Type != lexer.RPAREN {
		return nil, fmt.Errorf("expected ')' at position %d, got %q", p.curToken.Pos, p.curToken.Literal)
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return &ast.CallExpr{Function: function, Args: args}, nil
}
