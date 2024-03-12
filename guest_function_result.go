package wasify

import (
	"errors"
	"fmt"
	"io"

	"github.com/wasify-io/wasify-go/internal/types"
	"github.com/wasify-io/wasify-go/internal/utils"
)

type GuestFunctionResult interface {
	io.Closer
	Error() error
	Values() []PackedData

	ReadAnyPack(index int) (any, uint32, uint32, error)
	ReadBytesPack(index int) ([]byte, error)
	ReadBytePack(index int) (byte, error)
	ReadUint32Pack(index int) (uint32, error)
	ReadUint64Pack(index int) (uint64, error)
	ReadFloat32Pack(index int) (float32, error)
	ReadFloat64Pack(index int) (float64, error)
	ReadStringPack(index int) (string, error)
}

func NewGuestFunctionResult(err error, data uint64, memory Memory) GuestFunctionResult {
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
	memory Memory
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

func read(data uint64, memory Memory) ([]PackedData, error) {
	if data == 0 {
		return nil, errors.New("packedData is empty")
	}

	t, offsetU32, size := utils.UnpackUI64(data)

	if t != types.ValueTypePack {
		err := fmt.Errorf("Can't unpack host data, the type is not a valueTypePack. expected %d, got %d", types.ValueTypePack, t)
		return nil, err
	}

	bytes, err := memory.ReadBytes(offsetU32, size)
	if err != nil {
		err := errors.Join(errors.New("ReadPacks error, can't read bytes:"), err)
		return nil, err
	}

	err = memory.FreePack(PackedData(data))
	if err != nil {
		err := errors.Join(errors.New("ReadPacks error, can't free multiPackedData:"), err)
		return nil, err
	}

	return utils.Map(
		utils.BytesToUint64Array(bytes),
		func(data uint64) PackedData {
			return PackedData(data)
		},
	), nil
}
