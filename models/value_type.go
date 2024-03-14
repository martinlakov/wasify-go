package models

import "github.com/wasify-io/wasify-go/internal/types"

// ValueType represents the type of value used in function parameters and returns.
type ValueType types.ValueType

// supported value types in params and returns
const (
	ValueTypeBytes  = ValueType(types.ValueTypeBytes)
	ValueTypeByte   = ValueType(types.ValueTypeByte)
	ValueTypeI32    = ValueType(types.ValueTypeI32)
	ValueTypeI64    = ValueType(types.ValueTypeI64)
	ValueTypeF32    = ValueType(types.ValueTypeF32)
	ValueTypeF64    = ValueType(types.ValueTypeF64)
	ValueTypeString = ValueType(types.ValueTypeString)
)
