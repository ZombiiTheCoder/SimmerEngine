package vm

import (
	"fmt"
	"reflect"
	"simmer_js_engine/compiler/code"
	"simmer_js_engine/compiler/object"
)

func (vm *VM) ExecutePrefix(op code.Opcode) error {
	right := vm.pop()
	if reflect.TypeOf(right).Name() == "Number" {
		return vm.ExecutePrefixNumberOp(op, right.(object.Number))
	}
	return fmt.Errorf(
		"unsupported types for operation %v %v",
		code.Definitions[op].Name,
		reflect.TypeOf(right).Name(),
	)
}

func (vm *VM) ExecutePrefixNumberOp(op code.Opcode, right object.Number) error {
	r := right.Value
	var res float64
	switch op {
	case code.OP_UPLUS: res = r
	case code.OP_UMINUS: res = -r
	case code.OP_BNOT: res = float64(^int(r))
	default: return fmt.Errorf("unknown number op: %d", op)
	}
	return vm.push(object.Number{Value: res})
}

/*
case code.OP_PRINC:
case code.OP_PRDEC:	
*/