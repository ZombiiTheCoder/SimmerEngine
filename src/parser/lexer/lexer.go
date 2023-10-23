package lexer

import (
	"simmer_js_engine/parser/tokens"
	"strings"
)

type Lexer struct {
	ptr    int
	chars []byte
	tokens []tokens.Token
	isEof bool
	Line int
}

func InitLexer(chars []byte, _ error) *Lexer { return &Lexer{
	ptr: 0,
	chars: chars,
	tokens: make([]tokens.Token, 0),
	isEof: false,
} }

func isNum(char byte) bool { return char >= '0' && char <= '9' }
func isAlpha(char byte) bool { return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_' || char == '$' }
func isAlnum(char byte) bool { return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_' || char == '$' || char >= '0' && char <= '9' }

func (l *Lexer) At() byte { l.isEof = l.ptr >= len(l.chars); if !l.isEof { return l.chars[l.ptr] }; return ' '; }
func (l *Lexer) Peek(peek int) byte { isEof := l.ptr+peek >= len(l.chars); if !isEof { return l.chars[l.ptr+peek] }; return ' '; }
func (l *Lexer) Next() { l.ptr++ }

func (l *Lexer) addToken(value string, typee tokens.TokenType) {l.tokens = append(l.tokens, tokens.Token{
	Value: value,
	Type: typee,
	Line: l.Line,
}) }

func (l *Lexer) Tokenize() []tokens.Token {
	for !l.isEof {
		// Skip White Tokens
		switch l.At() { case ' ', '\t', '\b', '\r', '\f': l.Next(); case '\n': l.Next(); l.Line++ }
		// One/Multi CharTokens
		switch l.At() { 
		case '=':
			if l.Peek(1) == '=' { l.addToken("==", tokens.Equals); l.Next(); l.Next(); break }
			l.addToken("=", tokens.EqualsAssign); l.Next()
		case '<':
			if l.Peek(1) == '=' { l.addToken("<=", tokens.SmallerThanEqualTo); l.Next(); l.Next(); break }
			if l.Peek(1) == '<' { l.addToken("<<", tokens.LeftShift); l.Next(); l.Next(); break }
			if l.Peek(1) == '<' && l.Peek(2) == '=' { l.addToken("<<=", tokens.LeftShiftAssign); l.Next(); l.Next(); break }
			l.addToken("<", tokens.SmallerThan); l.Next()
		case '>':
			if l.Peek(1) == '=' { l.addToken(">=", tokens.BiggerThanEqualTo); l.Next(); l.Next(); break }
			if l.Peek(1) == '>' { l.addToken(">>", tokens.RightShift); l.Next(); l.Next(); break }
			if l.Peek(1) == '>'&& l.Peek(2) == '=' { l.addToken(">>=", tokens.RightShiftAssign); l.Next(); l.Next(); break }
			if l.Peek(1) == '>' && l.Peek(2) == '>' { l.addToken(">>>", tokens.UnsignedRightShift); l.Next(); l.Next(); l.Next(); break }
			if l.Peek(1) == '>' && l.Peek(2) == '>' && l.Peek(3) == '=' { l.addToken(">>>=", tokens.UnsignedRightShiftAssign); l.Next(); l.Next(); l.Next(); l.Next(); break }
			l.addToken(">", tokens.BiggerThan); l.Next()
		case '+':
			if l.Peek(1) == '+' { l.addToken("++", tokens.Increment); l.Next(); l.Next(); break }
			if l.Peek(1) == '=' { l.addToken("+=", tokens.PlusAssign); l.Next(); l.Next(); break }
			l.addToken("+", tokens.Plus); l.Next()
		case '-':
			if l.Peek(1) == '-' { l.addToken("--", tokens.Decrement); l.Next(); l.Next(); break }
			if l.Peek(1) == '=' { l.addToken("-=", tokens.MinusAssign); l.Next(); l.Next(); break }
			l.addToken("-", tokens.Minus); l.Next()
		case '!':
			if l.Peek(1) == '=' { l.addToken("!=", tokens.NotEquals); l.Next(); l.Next(); break }
			l.addToken("!", tokens.LogicalNot); l.Next()
		case '*':
			if l.Peek(1) == '=' { l.addToken("*=", tokens.MultiplyAssign); l.Next(); l.Next(); break }
			l.addToken("*", tokens.Multiply); l.Next()
		case '/':
			if l.Peek(1) == '=' { l.addToken("/=", tokens.DivideAssign); l.Next(); l.Next(); break }
			l.addToken("/", tokens.Divide); l.Next()
		case '&':
			if l.Peek(1) == '=' { l.addToken("&=", tokens.AndAssign); l.Next(); l.Next(); break }
			if l.Peek(1) == '&' { l.addToken("&&", tokens.LogicalAnd); l.Next(); l.Next(); break }
			l.addToken("&", tokens.BitwiseAnd); l.Next()
		case '|':
			if l.Peek(1) == '=' { l.addToken("|=", tokens.OrAssign); l.Next(); l.Next(); break }
			if l.Peek(1) == '|' { l.addToken("||", tokens.LogicalOr); l.Next(); l.Next(); break }
			l.addToken("|", tokens.BitwiseOr); l.Next()
		case '^':
			if l.Peek(1) == '=' { l.addToken("^=", tokens.XorAssign); l.Next(); l.Next(); break }
			l.addToken("^", tokens.BitwiseXor); l.Next()
		case '%':
			if l.Peek(1) == '=' { l.addToken("%=", tokens.ModuloAssign); l.Next(); l.Next(); break }
			l.addToken("%", tokens.Modulo); l.Next()
		case '~':
			if l.Peek(1) == '=' { l.addToken("~=", tokens.BitwiseNotAssign); l.Next(); l.Next(); break }
			l.addToken("~", tokens.BitwiseNot); l.Next()
		case ',': l.addToken(",", tokens.Comma); l.Next()
		case '.': l.addToken(".", tokens.Dot); l.Next()
		case '?': l.addToken("?", tokens.TernaryIf); l.Next()
		case ':': l.addToken(":", tokens.TernaryElse); l.Next()
		case ';': l.addToken(";", tokens.SemiColon); l.Next()
		case '(': l.addToken("(", tokens.OpenParen); l.Next()
		case ')': l.addToken(")", tokens.CloseParen); l.Next()
		case '{': l.addToken("{", tokens.OpenBrace); l.Next()
		case '}': l.addToken("}", tokens.CloseBrace); l.Next()
		}

		if isNum(l.At()) {
			value := new(strings.Builder)
			for isNum(l.At()) && !l.isEof {
				value.WriteByte(l.At()); l.Next()
			}
			if l.At() == '.' {
				value.WriteByte(l.At()); l.Next()
				for isNum(l.At()) && !l.isEof {
					value.WriteByte(l.At()); l.Next()
				}
			}
			l.addToken(value.String(), tokens.Number)
		}
		if isAlpha(l.At()) {
			value := new(strings.Builder)
			for isAlnum(l.At()) && !l.isEof {
				value.WriteByte(l.At()); l.Next()
			}
			l.addToken(value.String(), tokens.GetKeyword(value.String()))
		}
		// Tokenize String
		if l.At() == '"' || l.At() == '`' || l.At() == '\'' {
			l.Next()
			pos := l.ptr
			for l.At() == '"' || l.At() == '`' || l.At() == '\'' && !l.isEof { l.Next() }
			l.addToken(string(l.chars[pos:l.ptr]), tokens.String)
			l.Next()
		}
	}
	l.addToken("End Of File", tokens.EOF)
	return l.tokens;
}