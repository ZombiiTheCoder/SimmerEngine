package parser

import (
	"fmt"
	"os"
	"simmer_js_engine/parser/ast"
	"simmer_js_engine/parser/tokens"
)

type Parser struct {
	ptr    int
	tokens []tokens.Token
	isEof  bool
}

func InitParser(tokens []tokens.Token) *Parser {
	return &Parser{
		ptr:    0,
		tokens: tokens,
		isEof:  false,
	}
}

func (p *Parser) NotEof() bool { return p.At().Type != tokens.EOF }

func (p *Parser) At() tokens.Token { return p.tokens[p.ptr] }
func (p *Parser) Expect(ExpectedType tokens.TokenType) tokens.Token { p.ptr++; if p.tokens[p.ptr-1].Type == ExpectedType { return p.tokens[p.ptr-1] }; fmt.Println("Unexpected Type -> VALUE: "+p.tokens[p.ptr-1].Value); os.Exit(1); return tokens.Token{Value: "EOF", Type: tokens.EOF}; }
func (p *Parser) Next() tokens.Token { p.ptr++; return p.tokens[p.ptr-1] }

func (p *Parser) ProduceAst() ast.Stmt {
	body := make([]ast.Stmt, 0);
	for p.NotEof() { body = append(body, p.ParseStmt()) }
	return ast.Program{ Body: body }
}

func (p *Parser) ParseStmt() ast.Stmt {
	switch p.At().Type {
	default: return ast.Stmt(ast.ExprStmt{
		Expression: p.ParseExpr(),
	})
	}
}

func (p *Parser) ParseExpr() ast.Expr { return p.ParseAssignExpr() }

func (p *Parser) ParseAssignExpr() ast.Expr {
	left := p.ParseAdditiveExpr()
	if p.At().Type == tokens.EqualsAssign {
		op := p.Next().Type
		right := p.ParseAssignExpr()
		return ast.AssignExpr{
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseAdditiveExpr() ast.Expr {
	left := p.ParseMultiplicativeExpr()
	for p.At().Type == tokens.Plus || p.At().Type == tokens.Minus {
		op := p.Next().Type
		right := p.ParseMultiplicativeExpr()
		left = ast.BinaryExpr{
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseMultiplicativeExpr() ast.Expr {
	left := p.ParsePrimaryExpr()
	for p.At().Type == tokens.Multiply || p.At().Type == tokens.Divide || p.At().Type == tokens.Modulo {
		op := p.Next().Type
		right := p.ParsePrimaryExpr()
		left = ast.BinaryExpr{
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParsePrimaryExpr() ast.Expr {
	switch p.At().Type {
	case tokens.Number: return ast.NumberLiteral{
		Value: p.Next().Value,
	}
	default: fmt.Println("Unexpected Token: "+p.At().Value); os.Exit(-1); return ast.NumberLiteral{ Value: p.Next().Value }
	}
}