package tokens

type Token struct {
	Value string
	Type  TokenType
	Line  int
}

type TokenType string

const (
	Invalid TokenType = "!!!!!!!!!!!!!!!!!!"

	Equals       = "=="
	EqualsAssign = "="

	SmallerThanEqualTo = "<="
	LeftShift          = ">="
	LeftShiftAssign    = "<<="
	SmallerThan        = "<"

	BiggerThanEqualTo        = ">="
	RightShift               = ">>"
	RightShiftAssign         = ">>="
	UnsignedRightShift       = ">>>"
	UnsignedRightShiftAssign = ">>>="
	BiggerThan               = ">"

	Increment  = "++"
	PlusAssign = "+="
	Plus       = "+"

	Decrement   = "--"
	MinusAssign = "-="
	Minus       = "-"

	NotEquals  = "!="
	LogicalNot = "!"

	MultiplyAssign = "*="
	Multiply       = "*"

	DivideAssign = "/="
	Divide       = "/"

	AndAssign  = "&="
	LogicalAnd = "&&"
	BitwiseAnd = "&"

	OrAssign  = "|="
	LogicalOr = "||"
	BitwiseOr = "|"

	XorAssign  = "^="
	BitwiseXor = "^"

	ModuloAssign = "%="
	Modulo       = "%"

	BitwiseNotAssign = "~="
	BitwiseNot       = "~"

	Comma = ","
	Dot   = "."

	TernaryIf   = "?"
	TernaryElse = ":"

	SemiColon = ";"

	Number     = "Number"
	Identifier = "Identifier"
	String     = "String"

	OpenParen  = "("
	CloseParen = ")"

	OpenBrace  = "{"
	CloseBrace = "}"

	// Keywords
	Break    = "break"
	For      = "for"
	New      = "new"
	Var      = "var"
	Continue = "continue"
	Function = "function"
	Return   = "return"
	Void     = "void"
	Delete   = "delete"
	If       = "if"
	This     = "this"
	While    = "while"
	Else     = "else"
	In       = "in"
	Typeof   = "typeof"
	With     = "with"

	// Future Keywords
	Case     = "case"
	Debugger = "debugger"
	Export   = "export"
	Super    = "super"
	Catch    = "catch"
	Default  = "default"
	Extends  = "extends"
	Switch   = "switch"
	Class    = "class"
	Do       = "do"
	Finally  = "finally"
	Throw    = "throw"
	Const    = "const"
	Enum     = "enum"
	Import   = "import"
	Try      = "try"

	EOF = "EOF"
)

func GetKeyword(keyword string) TokenType {
	switch keyword {

	case "break":
		return Break
	case "for":
		return For
	case "new":
		return New
	case "var":
		return Var
	case "continue":
		return Continue
	case "function":
		return Function
	case "return":
		return Return
	case "void":
		return Void
	case "delete":
		return Delete
	case "if":
		return If
	case "this":
		return This
	case "while":
		return While
	case "else":
		return Else
	case "in":
		return In
	case "typeof":
		return Typeof
	case "with":
		return With
	case "case":
		return Case
	case "debugger":
		return Debugger
	case "export":
		return Export
	case "super":
		return Super
	case "catch":
		return Catch
	case "default":
		return Default
	case "extends":
		return Extends
	case "switch":
		return Switch
	case "class":
		return Class
	case "do":
		return Do
	case "finally":
		return Finally
	case "throw":
		return Throw
	case "const":
		return Const
	case "enum":
		return Enum
	case "import":
		return Import
	case "try":
		return Try
	default:
		return Identifier
	}
}