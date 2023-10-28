package parser

import (
	"fmt"
	"os"
	"simmer_js_engine/parser/ast"
	"simmer_js_engine/parser/tokens"
	"strconv"
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
func (p *Parser) Expect(ExpectedType tokens.TokenType) tokens.Token { p.ptr++; if p.tokens[p.ptr-1].Type == ExpectedType { return p.tokens[p.ptr-1] }; fmt.Println("Unexpected Type -> \""+string(p.tokens[p.ptr-1].Type)+"\" Expected -> \""+string(ExpectedType)+"\""+" Line: "+strconv.Itoa(p.At().Line+1)); os.Exit(1); return tokens.Token{Value: "EOF", Type: tokens.EOF}; }
func (p *Parser) ExpectMsg(ExpectedType tokens.TokenType, Msg string) tokens.Token { p.ptr++; if p.tokens[p.ptr-1].Type == ExpectedType { return p.tokens[p.ptr-1] }; fmt.Println("Unexpected Type -> \""+string(p.tokens[p.ptr-1].Type)+"\" Expected -> \""+string(ExpectedType)+"\""+" Line: "+strconv.Itoa(p.At().Line+1)); fmt.Println(Msg); os.Exit(1); return tokens.Token{Value: "EOF", Type: tokens.EOF}; }
func (p *Parser) Next() tokens.Token { p.ptr++; return p.tokens[p.ptr-1] }
func (p *Parser) Peek(pr int) tokens.Token { if (p.NotEof() && p.ptr+pr < len(p.tokens)) { return p.tokens[p.ptr+pr] }; return tokens.Token{Type: tokens.EOF} }
func (p *Parser) Past() tokens.Token { if (p.NotEof() && p.ptr-1 < len(p.tokens)) { return p.tokens[p.ptr-1] }; return tokens.Token{Type: tokens.EOF} }

func (p *Parser) ProduceAst() ast.Stmt {
	body := make([]ast.Stmt, 0);
	for p.NotEof() { body = append(body, p.ParseStmt()) }
	return ast.Program{ NodeType:"Program", Body: body }
}

func insert(a []tokens.Token, index int, value tokens.Token) []tokens.Token {
    if len(a) == index { // nil or empty slice or after last element
        return append(a, value)
    }
    a = append(a[:index+1], a[index:]...) // index < len(a)
    a[index] = value
    return a
}

func (p *Parser) ParseStmt() ast.Stmt {
	switch p.At().Type {
	case tokens.With: fmt.Println("WITH STATEMENT IS NOT IMPLEMENTED AND MAY NOT BE IMPLEMENTED"); os.Exit(-1); return ast.ContinueStmt{};
	case tokens.Continue: {
		p.Expect(tokens.Continue)
		if p.At().Type == tokens.SemiColon { p.Expect(tokens.SemiColon) }
		return ast.Stmt(ast.ContinueStmt{ NodeType:"ContinueStmt" })
	}
	case tokens.Break: {
		p.Expect(tokens.Break)
		if p.At().Type == tokens.SemiColon { p.Expect(tokens.SemiColon) }
		return ast.Stmt(ast.BreakStmt{ NodeType:"BreakStmt" })
	}
	case tokens.Return: return p.ParseReturnStmt()
	case tokens.While: return p.ParseWhileStmt()
	case tokens.If: return p.ParseIfStmt()
	case tokens.For: return p.ParseForStmt()
	case tokens.Var, tokens.Const: return p.ParseVarStmt()
	case tokens.Function: return p.ParseFuncStmt()
	default: {
		stmt := ast.Stmt(ast.ExprStmt{
			NodeType: "ExprStmt",
			Expression: p.ParseExpr(),
		})
		if p.At().Type == tokens.SemiColon { p.Expect(tokens.SemiColon) }
		return stmt;
	}
	}
}

func (p *Parser) ParseReturnStmt() ast.Expr {
	var expr ast.Expr
	p.Expect(tokens.Return)
	if p.At().Type == tokens.SemiColon { expr = ast.Identifier{
		NodeType: "Identifier",
		Value: "null",
	} } else {
		expr = p.ParseExpr()
		if p.At().Type == tokens.SemiColon { p.Expect(tokens.SemiColon) }
	}
	return ast.ReturnStmt{
		NodeType: "ReturnStmt",
		Expression: expr,
	}
}

func (p *Parser) ParseIfStmt() ast.Expr {
	var consequent ast.Expr
	var other ast.Expr
	line := p.Expect(tokens.If).Line
	p.Expect(tokens.OpenParen)
	condition := p.ParseExpr()
	p.Expect(tokens.CloseParen)
	if p.At().Type != tokens.OpenBrace { consequent = p.ParseStmt();
		if p.At().Type == tokens.Else && p.At().Line == line && p.Past().Type != tokens.SemiColon {
			p.ExpectMsg(tokens.SemiColon, "Expected \";\" or \"New Line\" for the \"If Statement\"")
		}	
	} else { consequent = p.ParseBodyStmt() }
	if p.At().Type == tokens.Else {
		p.Expect(tokens.Else)
		if p.At().Type == tokens.If { other = p.ParseIfStmt() } else {
			if p.At().Type != tokens.OpenBrace { other = p.ParseStmt() } else { other = p.ParseBodyStmt() }
		}
	}
	return ast.IfStmt{
		NodeType: "IfStmt",
		Condition: condition,
		Consequent: consequent,
		Other: other,
	}
}

func (p *Parser) ParseWhileStmt() ast.Expr {
	var body ast.Expr
	p.Expect(tokens.While)
	p.Expect(tokens.OpenParen)
	condition := p.ParseExpr()
	p.Expect(tokens.CloseParen)
	if p.At().Type != tokens.OpenBrace { body = p.ParseStmt(); } else { body = p.ParseBodyStmt() }
	return ast.WhileStmt{
		NodeType: "WhileStmt",
		Condition: condition,
		Body: body,
	}
}

func (p *Parser) ParseBodyStmt() ast.BodyStmt {
	p.Expect(tokens.OpenBrace)
	body := make([]ast.Stmt, 0);
	for p.NotEof() && p.At().Type != tokens.CloseBrace { body = append(body, p.ParseStmt()) }
	p.Expect(tokens.CloseBrace)
	if p.At().Type == tokens.SemiColon { p.Next() }
	return ast.BodyStmt{
		NodeType: "BodyStmt",
		Body: body,
	}
}

func (p *Parser) ParseArguments() []ast.Expr {
	p.Expect(tokens.OpenParen)
	if p.At().Type == tokens.CloseParen { p.Expect(tokens.CloseParen); return make([]ast.Expr, 0) }
	args := p.ParseArgList()
	p.Expect(tokens.CloseParen)
	return args
}

func (p *Parser) ParseArgList() []ast.Expr {
	args := make([]ast.Expr, 0)
	args = append(args, p.ParseAssignExpr())
	for (p.At().Type == tokens.Comma) {
		p.Next(); args = append(args, p.ParseAssignExpr())
	}
	return args
}

func (p *Parser) ParseFuncStmt() ast.Stmt {
	p.Expect(tokens.Function)
	name := p.Expect(tokens.Identifier).Value
	args := p.ParseArguments()
	body := p.ParseBodyStmt()
	
	return ast.FuncStmt{
		NodeType: "FuncStmt",
		Name: name,
		Args: args,
		Body: body.Body,
	}
}

func (p *Parser) ParseVarStmt() ast.Stmt {
	isConst := p.Next().Type == tokens.Const
	name := p.Expect(tokens.Identifier).Value
	p.Expect(tokens.EqualsAssign)
	expr := p.ParseExpr()
	if p.At().Type == tokens.Comma {
		return p.ParseVarStmtList_0(isConst, ast.VarStmt{
			NodeType: "VarStmt",
			Name: name,
			Value: expr,
			IsConst: isConst,
		})
	}
	if p.At().Type == tokens.SemiColon { p.Next() }
	return ast.VarStmt{
		NodeType: "VarStmt",
		Name: name,
		Value: expr,
		IsConst: isConst,
	}
}

func (p *Parser) ParseVarStmtList_0(isConst bool, init ast.VarStmt) ast.Stmt {
	decs := make([]ast.Stmt, 0); decs = append(decs, init)
	for p.At().Type == tokens.Comma {
		p.Expect(tokens.Comma); decs = append(decs, p.ParseSimpleVarStmt(isConst))
	}
	if p.At().Type == tokens.SemiColon { p.Next() }
	return ast.VarStmtList{
		NodeType: "VarStmtList",
		Stmts: decs,
	}
}

func (p *Parser) ParseSimpleVarStmt(isConst bool) ast.Stmt {
	name := p.Expect(tokens.Identifier).Value
	expr := ast.Expr(ast.Identifier{
		NodeType: "Identifier",
		Value: "null",
	})
	if p.At().Type == tokens.EqualsAssign {
		p.Next()
		expr = p.ParseExpr()
	}
	return ast.VarStmt{
		NodeType: "VarStmt",
		Name: name,
		Value: expr,
		IsConst: isConst,
	}
}

func (p *Parser) ParseForStmt() ast.Stmt {
	var init ast.StmtExpr
	var condition ast.Expr
	var after ast.Expr
	var stmt ast.Stmt

	p.Expect(tokens.For)
	p.Expect(tokens.OpenParen)
	if p.Peek(2).Type == tokens.In { return p.ParseForInStmt() }
	if p.At().Type != tokens.SemiColon {
		if (p.At().Type == tokens.Var || p.At().Type == tokens.Const) {
			init = p.ParseVarStmt()
		} else {
			init = p.ParseExpr();
		}
		if p.Past().Type != tokens.SemiColon { p.Expect(tokens.SemiColon) }
	} else { p.Expect(tokens.SemiColon) } 
	if p.At().Type != tokens.SemiColon {
		condition = p.ParseExpr(); p.Expect(tokens.SemiColon)
	} else { p.Expect(tokens.SemiColon) } 
	if p.At().Type != tokens.CloseParen {
		after = p.ParseExpr()
		p.Expect(tokens.CloseParen)
	} else { p.Expect(tokens.CloseParen) } 
	if p.At().Type != tokens.OpenBrace {
		stmt = p.ParseStmt()
	} else {
		stmt = p.ParseBodyStmt()
	}
	
	return ast.ForStmt{
		NodeType: "ForStmt",
		Init: init,
		Condition: condition,
		After: after,
		Statement: stmt,
	}
}

func (p *Parser) ParseForInStmt() ast.Stmt {
	var stmt ast.Stmt
	isConst := p.Next().Type == tokens.Const
	name := p.Next().Value
	p.Expect(tokens.In)
	object := p.ParseExpr()
	p.Expect(tokens.CloseParen)
	if p.At().Type != tokens.OpenBrace {
		stmt = p.ParseStmt()
	} else {
		stmt = p.ParseBodyStmt()
	}

	return ast.ForInStmt{
		NodeType: "ForInStmt",
		Var: name,
		IsConst: isConst,
		Object: object,
		Statement: stmt,
	}
}


func (p *Parser) ParseExpr() ast.Expr { return p.ParseAssignExpr() }

func (p *Parser) ParseAssignExpr() ast.Expr {
	left := p.ParseTernaryExpr()
	switch t := p.At().Type; t {
	case tokens.EqualsAssign, tokens.PlusAssign, tokens.MinusAssign, tokens.MultiplyAssign,
		tokens.DivideAssign, tokens.ModuloAssign, tokens.LeftShiftAssign, tokens.RightShiftAssign,
		tokens.UnsignedRightShiftAssign, tokens.AndAssign, tokens.XorAssign, tokens.OrAssign, tokens.BitwiseNotAssign:
		return ast.AssignExpr{
			NodeType: "AssignExpr",
			Left:  left,
			Op:    p.Next().Type,
			Right: p.ParseAssignExpr(),
		}
	}
	return left
}

func (p *Parser) ParseTernaryExpr() ast.Expr {
	left := p.ParseLogicalOrExpr()
	if p.At().Type == tokens.TernaryIf {
		p.Expect(tokens.TernaryIf)
		iif := p.ParseAssignExpr()
		p.Expect(tokens.TernaryElse)
		other := p.ParseAssignExpr()
		return ast.TernaryExpr{
			NodeType: "TernaryExpr",
			Condition: left,
			Consequent: iif,
			Other: other,
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
	for p.At().Type == tokens.Equals || p.At().Type == tokens.NotEquals {
		op := p.Next().Type
		right := p.ParseRelationalExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) ParseRelationalExpr() ast.Expr {
	left := p.ParseShiftExpr()
	for p.At().Type == tokens.SmallerThan || p.At().Type == tokens.SmallerThanEqualTo ||
		p.At().Type == tokens.BiggerThan || p.At().Type == tokens.BiggerThanEqualTo ||
		p.At().Type == tokens.In {
		op := p.Next().Type
		right := p.ParseShiftExpr()
		left = ast.BinaryExpr{
			NodeType: "BinaryExpr",
			Left:  left,
			Op:    op,
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
	if p.At().Type == tokens.Increment ||
	p.At().Type == tokens.Decrement ||
	p.At().Type == tokens.LogicalNot ||
	p.At().Type == tokens.BitwiseNot ||
	p.At().Type == tokens.Plus ||
	p.At().Type == tokens.Minus ||
	p.At().Type == tokens.Typeof ||
	p.At().Type == tokens.Void ||
	p.At().Type == tokens.Delete {
		op := p.Next().Type
		right := p.ParseAssignExpr()
		return ast.PrefixExpr{
			NodeType: "PrefixExpr",
			Op:   op,
			Right: right,
		}
	}

	return p.ParsePostfixExpr()
}

func (p *Parser) ParsePostfixExpr() ast.Expr {
	left := p.ParseCallMemberExpr()
	if left.GetNodeType() != "Identifier" { return left }
	if p.At().Type == tokens.Increment || p.At().Type == tokens.Decrement {
		op := p.Next().Type
		return ast.PostfixExpr{
			NodeType: "PostfixExpr",
			Op: op,
			Left: left,
		}
	}

	return left

}

func (p *Parser) ParseCallMemberExpr() ast.Expr {
	member := p.ParsePrimaryExpr()
	if p.At().Type == tokens.OpenParen && member.GetNodeType() == "Identifier" && p.Past().Line == p.At().Line { return p.ParseCallExpr(member) }
	return member
}

func (p *Parser) ParseCallExpr(caller ast.Expr) ast.Expr {
	callExpr := ast.CallExpr{
		Caller: caller,
		Args: p.ParseArguments(),
	}
	if p.At().Type == tokens.OpenParen { callExpr = p.ParseCallExpr(callExpr).(ast.CallExpr) }
	return callExpr
}

func (p *Parser) ParsePrimaryExpr() ast.Expr {
	switch p.At().Type {
	case tokens.Identifier:
		return ast.Identifier{
			NodeType: "Identifier",
			Value:    p.Next().Value,
	}
	case tokens.Number:
		return ast.Literal{
			NodeType: "NumberLiteral",
			Value:    p.Next().Value,
	}
	case tokens.String:
		return ast.Literal{
			NodeType: "StringLiteral",
			Value:    p.Next().Value,
	}
	case tokens.True, tokens.False:
		return ast.Literal{
			NodeType: "BooleanLiteral",
			Value:    p.Next().Value,
	}
	case tokens.Null:
		return ast.Literal{
			NodeType: "NullLiteral",
			Value:    p.Next().Value,
	}
	case tokens.OpenParen:
		p.Expect(tokens.OpenParen)
		expr := p.ParseExpr()
		p.Expect(tokens.CloseParen)
		return expr
	default:
		fmt.Println("Unexpected Token: " + fmt.Sprint(p.At().Type) +" Line: "+strconv.Itoa(p.At().Line+1))
		os.Exit(-1)
		return ast.Literal{Value: p.At().Value}
	}
}