package types

import (
	"fmt"
	"reflect"
)

// ValueType is an enumeration of supported data types for function parameters and returns.
type ValueType uint8

// ValueTypePack is a reserved ValueType used for packed data.
const ValueTypePack ValueType = 255

// These constants represent the possible data types that can be used in function parameters and returns.
const (
	ValueTypeBytes ValueType = iota
	ValueTypeByte
	ValueTypeI32
	ValueTypeI64
	ValueTypeF32
	ValueTypeF64
	ValueTypeString
)

func (v ValueType) String() string {
	switch v {
	case ValueTypePack:
		return "ValueTypePack"
	case ValueTypeBytes:
		return "ValueTypeBytes"
	case ValueTypeByte:
		return "ValueTypeByte"
	case ValueTypeI32:
		return "ValueTypeI32"
	case ValueTypeI64:
		return "ValueTypeI64"
	case ValueTypeF32:
		return "ValueTypeF32"
	case ValueTypeF64:
		return "ValueTypeF64"
	case ValueTypeString:
		return "ValueTypeString"
	}

	return "undefined"
}

// GetOffsetSizeAndDataTypeByConversion determines the memory size (offsetSize) and ValueType
// of a given data. The function supports several data
func GetOffsetSizeAndDataTypeByConversion(data any) (dataType ValueType, offsetSize uint32, err error) {
	switch vTyped := data.(type) {
	case []byte:
		return ValueTypeBytes, uint32(len(vTyped)), nil
	case byte:
		return ValueTypeByte, 1, nil
	case uint32:
		return ValueTypeI32, 4, nil
	case uint64:
		return ValueTypeI64, 8, nil
	case float32:
		return ValueTypeF32, 4, nil
	case float64:
		return ValueTypeF64, 8, nil
	case string:
		return ValueTypeString, uint32(len(vTyped)), nil
	default:
		return 0, 0, fmt.Errorf("unsupported conversion data type %s", reflect.TypeOf(vTyped))
	}
}
