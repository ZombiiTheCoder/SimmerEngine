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

	UnaryNotAssign

	Comma
	Dot

	UnaryIf
	UnaryElse

	SemiColon

	Number
	Identifer
	String

	OpenParen
	CloseParen

	OpenBrace
	CloseBrace

	EOF
)