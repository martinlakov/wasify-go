package wazero

import (
	"errors"
	"fmt"

	"github.com/wasify-io/wasify-go/internal/types"
	. "github.com/wasify-io/wasify-go/internal/utils"
	. "github.com/wasify-io/wasify-go/models"
)

func _NewGuestFunctionResult(err error, data MultiPackedData, memory RMemory) GuestFunctionResult {
	if err != nil {
		return &_GuestFunctionResult{err: err}
	}

	if types.ValueType(data>>56) != types.ValueTypePack {
		return &_GuestFunctionResult{memory: memory}
	}

	values, err := read(data, memory)
	if err != nil {
		return &_GuestFunctionResult{err: err}
	}

	return &_GuestFunctionResult{
		memory: memory,
		values: values,
	}
}

type _GuestFunctionResult struct {
	err    error
	memory RMemory
	values []PackedData
}

func (self *_GuestFunctionResult) Close() error {
	return self.memory.FreePack(self.values...)
}

func (self *_GuestFunctionResult) Error() error {
	return self.err
}

func (self *_GuestFunctionResult) Values() []PackedData {
	return self.values
}

func (self *_GuestFunctionResult) ReadAnyPack(index int) (any, uint32, uint32, error) {
	return self.memory.ReadAnyPack(self.values[index])
}

func (self *_GuestFunctionResult) ReadBytesPack(index int) ([]byte, error) {
	return self.memory.ReadBytesPack(self.values[index])
}

func (self *_GuestFunctionResult) ReadBytePack(index int) (byte, error) {
	return self.memory.ReadBytePack(self.values[index])
}

func (self *_GuestFunctionResult) ReadUint32Pack(index int) (uint32, error) {
	return self.memory.ReadUint32Pack(self.values[index])
}

func (self *_GuestFunctionResult) ReadUint64Pack(index int) (uint64, error) {
	return self.memory.ReadUint64Pack(self.values[index])
}

func (self *_GuestFunctionResult) ReadFloat32Pack(index int) (float32, error) {
	return self.memory.ReadFloat32Pack(self.values[index])
}

func (self *_GuestFunctionResult) ReadFloat64Pack(index int) (float64, error) {
	return self.memory.ReadFloat64Pack(self.values[index])
}

func (self *_GuestFunctionResult) ReadStringPack(index int) (string, error) {
	return self.memory.ReadStringPack(self.values[index])
}

func read(data MultiPackedData, memory RMemory) ([]PackedData, error) {
	if data == 0 {
		return nil, errors.New("packedData is empty")
	}

	t, offsetU32, size := UnpackUI64[ValueType](uint64(data))
	if t != ValueType(255) {
		return nil, fmt.Errorf("invalid data type found, expected %d, got %d", types.ValueTypePack, t)
	}

	bytes, err := memory.ReadBytes(offsetU32, size)
	if err != nil {
		return nil, errors.Join(errors.New("failed to read data"), err)
	}

	err = memory.FreePack(PackedData(data))
	if err != nil {
		return nil, errors.Join(errors.New("failed to free up pack data"), err)
	}

	return Map(
		BytesToUint64Array(bytes),
		func(data uint64) PackedData {
			return PackedData(data)
		},
	), nil
}
