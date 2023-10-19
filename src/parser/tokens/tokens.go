package tokens

type Token struct {
	Value string
	Type  TokenType
}

type TokenType int

const (
	Invalid TokenType = iota

	Equals
	EqualsAssign

	SmallerThanEqualTo
	LeftShift
	LeftShiftAssign
	SmallerThan

	BiggerThanEqualTo
	RightShift
	RightShiftAssign
	UnsignedRightShift
	UnsignedRightShiftAssign
	BiggerThan

	Increment
	PlusAssign
	Plus

	Decrement
	MinusAssign
	Minus

	NotEquals
	LogicalNot

	MultiplyAssign
	Multiply

	DivideAssign
	Divide

	AndAssign
	LogicalAnd
	BitwiseAnd

	OrAssign
	LogicalOr
	BitwiseOr

	XorAssign
	BitwiseXor

	ModuloAssign
	Modulo

	BitwiseNotAssign
	BitwiseNot

	Comma
	Dot

	TernaryIf
	TernaryElse

	SemiColon

	Number
	Identifer
	String

	OpenParen
	CloseParen

	OpenBrace
	CloseBrace

	// Keywords
	Break
	For
	New
	Var
	Continue
	Function
	Return
	Void
	Delete
	If
	This
	While
	Else
	In
	Typeof
	With

	// Future Keywords
	Case
	Debugger
	Export
	Super
	Catch
	Default
	Extends
	Switch
	Class
	Do
	Finally
	Throw
	Const
	Enum
	Import
	Try

	EOF
)

func GetKeyword(keyword string) TokenType {
	dict := map[string]TokenType{
		"break":    Break,
		"for":      For,
		"new":      New,
		"var":      Var,
		"continue": Continue,
		"function": Function,
		"return":   Return,
		"void":     Void,
		"delete":   Delete,
		"if":       If,
		"this":     This,
		"while":    While,
		"else":     Else,
		"in":       In,
		"typeof":   Typeof,
		"with":     With,

		"case":     Case,
		"debugger": Debugger,
		"export":   Export,
		"super":    Super,
		"catch":    Catch,
		"default":  Default,
		"extends":  Extends,
		"switch":   Switch,
		"class":    Class,
		"do":       Do,
		"finally":  Finally,
		"throw":    Throw,
		"const":    Const,
		"enum":     Enum,
		"import":   Import,
		"try":      Try,
	}
	if _, ok := dict[keyword]; ok {
		return dict[keyword]
	} else {
		return Identifer
	}
}