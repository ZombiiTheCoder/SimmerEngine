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
}) }

func (l *Lexer) Tokenize() []tokens.Token {
	for !l.isEof {
		// Skip White Tokens
		switch l.At() { case ' ', '\n', '\t', '\b', '\r', '\f': l.Next() }
		// One/Multi CharTokens
		switch l.At() { 
		case '=': l.Next()
			if l.Peek(1) == '=' { l.addToken("==", tokens.Equals); l.Next(); break }
			l.addToken("=", tokens.EqualsAssign)
		case '<': l.Next()
			if l.Peek(1) == '=' { l.addToken("<=", tokens.SmallerThanEqualTo); l.Next(); break }
			if l.Peek(1) == '<' { l.addToken("<<", tokens.LeftShift); l.Next(); break }
			if l.Peek(1) == '<' { l.addToken("<<=", tokens.LeftShiftAssign); l.Next(); break }
			l.addToken("<", tokens.SmallerThan)
		case '>': l.Next()
			if l.Peek(1) == '=' { l.addToken(">=", tokens.BiggerThanEqualTo); l.Next(); break }
			if l.Peek(1) == '>' { l.addToken(">>", tokens.RightShift); l.Next(); break }
			if l.Peek(1) == '>' { l.addToken(">>=", tokens.RightShiftAssign); l.Next(); break }
			if l.Peek(1) == '>' && l.Peek(2) == '>' { l.addToken(">>>", tokens.UnsignedRightShift); l.Next(); l.Next(); break }
			if l.Peek(1) == '>' && l.Peek(2) == '>' && l.Peek(3) == '=' { l.addToken(">>>=", tokens.UnsignedRightShiftAssign); l.Next(); l.Next(); l.Next(); break }
			l.addToken(">", tokens.BiggerThan)
		case '+': l.Next()
			if l.Peek(1) == '+' { l.addToken("++", tokens.Increment); l.Next(); break }
			if l.Peek(1) == '=' { l.addToken("+=", tokens.PlusAssign); l.Next(); break }
			l.addToken("+", tokens.Plus);
		case '-': l.Next()
			if l.Peek(1) == '-' { l.addToken("--", tokens.Decrement); l.Next(); break }
			if l.Peek(1) == '=' { l.addToken("-=", tokens.MinusAssign); l.Next(); break }
			l.addToken("-", tokens.Minus)
		case '!': l.Next()
			if l.Peek(1) == '=' { l.addToken("!=", tokens.NotEquals); l.Next(); break }
			l.addToken("!", tokens.LogicalNot)
		case '*': l.Next()
			if l.Peek(1) == '=' { l.addToken("*=", tokens.MultiplyAssign); l.Next(); break }
			l.addToken("*", tokens.Multiply)
		case '/': l.Next()
			if l.Peek(1) == '=' { l.addToken("/=", tokens.DivideAssign); l.Next(); break }
			l.addToken("/", tokens.Divide)
		case '&': l.Next()
			if l.Peek(1) == '=' { l.addToken("&=", tokens.AndAssign); l.Next(); break }
			if l.Peek(1) == '|' { l.addToken("&&", tokens.LogicalAnd); l.Next(); break }
			l.addToken("&", tokens.BitwiseAnd)
		case '|': l.Next()
			if l.Peek(1) == '=' { l.addToken("|=", tokens.OrAssign); l.Next(); break }
			if l.Peek(1) == '|' { l.addToken("||", tokens.LogicalOr); l.Next(); break }
			l.addToken("|", tokens.BitwiseOr)
		case '^': l.Next()
			if l.Peek(1) == '=' { l.addToken("^=", tokens.XorAssign); l.Next(); break }
			l.addToken("^", tokens.BitwiseXor);
		case '%': l.Next()
			if l.Peek(1) == '=' { l.addToken("%=", tokens.ModuloAssign); l.Next(); break }
			l.addToken("%", tokens.Modulo)
		case '~': l.addToken("~", tokens.UnaryNotAssign); l.Next()
		case ',': l.addToken(",", tokens.Comma); l.Next()
		case '.': l.addToken(".", tokens.Dot); l.Next()
		case '?': l.addToken("?", tokens.UnaryIf); l.Next()
		case ':': l.addToken(":", tokens.UnaryElse); l.Next()
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
			l.addToken(value.String(), tokens.Identifer)
		}
		// Tokenize String
		if l.At() == '"' || l.At() == '`' || l.At() == '\'' {
			l.Next()
			value := new(strings.Builder)
			for l.At() == '"' || l.At() == '`' || l.At() == '\'' && !l.isEof {
				value.WriteByte(l.At()); l.Next()
			}
			l.Next()
			l.addToken(value.String(), tokens.Identifer)
		}
	}
	l.addToken("End Of File", tokens.EOF)
	return l.tokens;
}