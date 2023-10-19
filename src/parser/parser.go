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
func (p *Parser) Peek() tokens.Token { if (p.NotEof() && p.ptr+1 < len(p.tokens)) { return p.tokens[p.ptr+1] }; return tokens.Token{Type: tokens.EOF} }

func (p *Parser) ProduceAst() ast.Stmt {
	body := make([]ast.Stmt, 0);
	for p.NotEof() { body = append(body, p.ParseStmt()) }
	return ast.Program{ NodeType:"Program", Body: body }
}

func (p *Parser) ParseStmt() ast.Stmt {
	switch p.At().Type {
	default: return ast.Stmt(ast.ExprStmt{
		NodeType: "ExprStmt",
		Expression: p.ParseExpr(),
	})
	}
}

func (p *Parser) ParseExpr() ast.Expr { return p.ParseAssignExpr() }

func (p *Parser) ParseAssignExpr() ast.Expr {
	left := p.ParseTernaryExpr()
	switch p.At().Type {
		case tokens.EqualsAssign,
		tokens.PlusAssign,
		tokens.MinusAssign,
		tokens.MultiplyAssign,
		tokens.DivideAssign,
		tokens.ModuloAssign,
		tokens.LeftShiftAssign,
		tokens.RightShiftAssign,
		tokens.UnsignedRightShiftAssign,
		tokens.AndAssign,
		tokens.XorAssign,
		tokens.OrAssign,
		tokens.BitwiseNotAssign:
			op := p.Next().Type
			right := p.ParseAssignExpr()
			return ast.AssignExpr{
				NodeType: `AssignExpr`,
				Left: left,
				Op: op,
				Right: right,
			}
	}
	return left
}

func (p *Parser) ParseTernaryExpr() ast.Expr {
	left := p.ParseAdditiveExpr()
	if p.At().Type == tokens.TernaryIf {
		p.Expect(tokens.TernaryIf)
		iif := p.ParseAssignExpr()
		p.Expect(tokens.TernaryElse)
		elsee := p.ParseAssignExpr()
		return ast.TernaryExpr{
			NodeType: "TernaryExpr",
			Condition: left,
			Consequent: iif,
			Else: elsee,
		}
	}
	return left
}

func (p *Parser) ParseLogicalOrExpr() ast.Expr {
	left := p.ParseLogicalAndExpr()
	for p.At().Type == tokens.LogicalOr {
		op := p.Next().Type
		right := p.ParseLogicalAndExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseLogicalAndExpr() ast.Expr {
	left := p.ParseBitwiseOrExpr()
	for p.At().Type == tokens.LogicalAnd {
		op := p.Next().Type
		right := p.ParseBitwiseOrExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseBitwiseOrExpr() ast.Expr {
	left := p.ParseBitwiseXorExpr()
	for p.At().Type == tokens.BitwiseOr {
		op := p.Next().Type
		right := p.ParseBitwiseXorExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}


func (p *Parser) ParseBitwiseXorExpr() ast.Expr {
	left := p.ParseBitwiseAndExpr()
	for p.At().Type == tokens.BitwiseXor {
		op := p.Next().Type
		right := p.ParseBitwiseAndExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}


func (p *Parser) ParseBitwiseAndExpr() ast.Expr {
	left := p.ParseEqualityExpr()
	for p.At().Type == tokens.BitwiseAnd {
		op := p.Next().Type
		right := p.ParseEqualityExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseEqualityExpr() ast.Expr {
	left := p.ParseRelationalExpr()
	for p.At().Type == tokens.Equals ||
	p.At().Type == tokens.NotEquals {
		op := p.Next().Type
		right := p.ParseRelationalExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseRelationalExpr() ast.Expr {
	left := p.ParseShiftExpr()
	for p.At().Type == tokens.SmallerThan ||
	p.At().Type == tokens.SmallerThanEqualTo ||
	p.At().Type == tokens.BiggerThan ||
	p.At().Type == tokens.BiggerThanEqualTo ||
	p.At().Type == tokens.In {
		op := p.Next().Type
		right := p.ParseShiftExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}


func (p *Parser) ParseShiftExpr() ast.Expr {
	left := p.ParseAdditiveExpr()
	for p.At().Type == tokens.LeftShift ||
	p.At().Type == tokens.RightShift ||
	p.At().Type == tokens.UnsignedRightShift {
		op := p.Next().Type
		right := p.ParseAdditiveExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
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
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseMultiplicativeExpr() ast.Expr {
	left := p.ParsePrefixExpr()
	for p.At().Type == tokens.Multiply || p.At().Type == tokens.Divide || p.At().Type == tokens.Modulo {
		op := p.Next().Type
		right := p.ParsePrefixExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left: left,
			Op: op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParsePrefixExpr() ast.Expr {
	fmt.Println(p.Peek().Type == tokens.Increment)
	if (p.Peek().Type == tokens.Increment) {
		op := p.Next().Type
		right := p.ParsePostfixExpr()
		return ast.PrefixExpr{
			NodeType: "PrefixExpr",
			Op: op,
			Right: right,
		}
	}

	return p.ParsePostfixExpr()
}

func (p *Parser) ParsePostfixExpr() ast.Expr {
	left := p.ParsePrimaryExpr()
	for p.At().Type == tokens.Increment || p.At().Type == tokens.Decrement {
		fmt.Println(p.At())
		op := p.Next().Type
		return ast.PostfixExpr{
			NodeType: "PostfixExpr",
			Left: left,
			Op: op,
		}
	}
	return left
}

func (p *Parser) ParsePrimaryExpr() ast.Expr {
	switch p.At().Type {
	case tokens.Identifer: return ast.Identifer{
		NodeType: "Identifer",
		Value: p.Next().Value,
	}
	case tokens.Number: return ast.NumberLiteral{
		NodeType: "NumberLiteral",
		Value: p.Next().Value,
	}
	case tokens.String: return ast.StringLiteral{
		NodeType: "StringLiteral",
		Value: p.Next().Value,
	}
	default: fmt.Println("Unexpected Token: "+p.At().Value); os.Exit(-1); return ast.NumberLiteral{ Value: p.Next().Value }
	}
}