package vm

import (
	"fmt"
	"math"
	"reflect"
	"simmer_js_engine/compiler/code"
	"simmer_js_engine/compiler/object"
)

func (vm *VM) ExecuteBinaryOp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if reflect.TypeOf(left).Name() == "Number" && reflect.TypeOf(right).Name() == "Number" {
		return vm.ExecuteBinaryNumberOp(op, left.(object.Number), right.(object.Number))
	}
	return fmt.Errorf(
		"unsupported types for operation %v %v %v",
		reflect.TypeOf(left).Name(),
		code.Definitions[op].Name,
		reflect.TypeOf(right).Name(),
	)
}

func (vm *VM) ExecuteBinaryNumberOp(op code.Opcode, right, left object.Number) error {
	l := left.Value
	r := right.Value
	var res float64
	switch op {
	case code.OP_ADD: res = l + r
	case code.OP_SUB: res = l - r
	case code.OP_MUL: res = l * r
	case code.OP_DIV: res = l / r
	// case code.OP_RSHIFT: res = ,math.
	case code.OP_MOD: res = math.Mod(l, r)
	default: return fmt.Errorf("unknown number op: %d", op)
	}
	return vm.push(object.Number{Value: res})
}

func (vm *VM) ExecuteComparison(op code.Opcode) error {
	left := vm.pop()
	right := vm.pop()
	if reflect.TypeOf(left).Name() == "Number" || reflect.TypeOf(right).Name() == "Number" {
		return vm.ExecuteComparisonNumberOp(op, left.(object.Number), right.(object.Number))
	}
	switch op {
	case code.OP_EQUAL: return vm.push(object.Boolean{Value: right == left})
	case code.OP_NOT_EQUAL: return vm.push(object.Boolean{Value: right != left})
	default:
		return fmt.Errorf(
			"unsupported types for operation %v %v %v",
			reflect.TypeOf(left).Name(),
			code.Definitions[op].Name,
			reflect.TypeOf(right).Name(),
		)
	}
}

func (vm *VM) ExecuteComparisonNumberOp(op code.Opcode, right, left object.Number) error {
	l := left.Value
	r := right.Value
	var res bool
	switch op {
	case code.OP_EQUAL: res = right == left
	case code.OP_NOT_EQUAL: res = right != left
	case code.OP_GREATER_THAN: res = l > r
	case code.OP_GREATER_THAN_EQUAL: res = l >= r
	case code.OP_LESS_THAN: res = l < r
	case code.OP_LESS_THAN_EQUAL: res = l <= r

	default: return fmt.Errorf("unknown number op: %d", op)
	}
	return vm.push(object.Boolean{Value: res})
}

