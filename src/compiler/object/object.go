package object

type Object interface {
	GetValue() any
}

type Number struct {
	Object
	Value float64
}

func (n Number) GetValue() any { return n.Value }

type Boolean struct {
	Object
	Value bool
}

func (n Boolean) GetValue() any { return n.Value }

type Null struct {
	Object
	Value interface{}
}

func (n Null) GetValue() any { return n.Value }

func IsTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case Boolean:
		return obj.Value
	case Number:
		return obj.Value != 0
	default:
		return true
	}
}