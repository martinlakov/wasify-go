package wazero

import (
	"fmt"
	"reflect"

	"github.com/wasify-io/wasify-go/internal/types"
	. "github.com/wasify-io/wasify-go/internal/utils"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

type _Memory struct {
	module *_Module
	logger logging.Logger
}

// ReadAnyPack extracts and reads data from a packed memory location.
//
// Given a packed data representation, this function determines the type, offset, and size of the data to be read.
// It then reads the data from the specified offset and returns it.
//
// Returns:
// - offset: The memory location where the data starts.
// - size: The size or length of the data.
// - data: The actual extracted data of the determined type (i.e., byte slice, uint32, uint64, float32, float64).
// - error: An error if encountered (e.g., unsupported data type, out-of-range error).
func (self *_Memory) ReadAnyPack(pd PackedData) (any, uint32, uint32, error) {
	var err error
	var data any

	// Unpack the packedData to extract offset and size values.
	valueType, offset, size := UnpackUI64[ValueType](uint64(pd))

	switch valueType {
	case ValueTypeBytes:
		data, err = self.ReadBytes(offset, size)
	case ValueTypeByte:
		data, err = self.ReadByte(offset)
	case ValueTypeI32:
		data, err = self.ReadUint32(offset)
	case ValueTypeI64:
		data, err = self.ReadUint64(offset)
	case ValueTypeF32:
		data, err = self.ReadFloat32(offset)
	case ValueTypeF64:
		data, err = self.ReadFloat64(offset)
	case ValueTypeString:
		data, err = self.ReadString(offset, size)
	default:
		err = fmt.Errorf("unsupported read data type %v", valueType)
	}

	if err != nil {
		self.logger.Error(err.Error())
		return nil, 0, 0, err
	}

	return data, offset, size, err
}

func (self *_Memory) ReadBytes(offset uint32, size uint32) ([]byte, error) {
	buf, ok := self.module.raw.Memory().Read(offset, size)
	if !ok {
		err := fmt.Errorf("Memory.ReadBytes(%d, %d) out of range of memory size %d", offset, size, self.Size())
		self.logger.Error(err.Error())
		return nil, err
	}

	return buf, nil
}

func (self *_Memory) ReadBytesPack(pd PackedData) ([]byte, error) {
	_, offset, size := UnpackUI64[ValueType](uint64(pd))
	return self.ReadBytes(offset, size)
}

func (self *_Memory) ReadByte(offset uint32) (byte, error) {
	buf, ok := self.module.raw.Memory().ReadByte(offset)
	if !ok {
		err := fmt.Errorf("Memory.ReadByte(%d, %d) out of range of memory size %d", offset, 1, self.Size())
		self.logger.Error(err.Error())
		return 0, err
	}

	return buf, nil
}

func (self *_Memory) ReadBytePack(pd PackedData) (byte, error) {
	return self.ReadByte(Second(UnpackUI64[ValueType](uint64(pd))))
}

func (self *_Memory) ReadUint32(offset uint32) (uint32, error) {
	data, ok := self.module.raw.Memory().ReadUint32Le(offset)
	if !ok {
		err := fmt.Errorf("Memory.ReadUint32(%d, %d) out of range of memory size %d", offset, 4, self.Size())
		self.logger.Error(err.Error())
		return 0, err
	}

	return data, nil
}

func (self *_Memory) ReadUint32Pack(pd PackedData) (uint32, error) {
	return self.ReadUint32(Second(UnpackUI64[ValueType](uint64(pd))))
}

func (self *_Memory) ReadUint64(offset uint32) (uint64, error) {
	data, ok := self.module.raw.Memory().ReadUint64Le(offset)
	if !ok {
		err := fmt.Errorf("Memory.ReadUint64(%d, %d) out of range of memory size %d", offset, 8, self.Size())
		self.logger.Error(err.Error())
		return 0, err
	}

	return data, nil
}

func (self *_Memory) ReadUint64Pack(pd PackedData) (uint64, error) {
	return self.ReadUint64(Second(UnpackUI64[ValueType](uint64(pd))))
}

func (self *_Memory) ReadFloat32(offset uint32) (float32, error) {
	data, ok := self.module.raw.Memory().ReadFloat32Le(offset)
	if !ok {
		err := fmt.Errorf("Memory.ReadFloat32(%d, %d) out of range of memory size %d", offset, 4, self.Size())
		self.logger.Error(err.Error())
		return 0, err
	}

	return data, nil
}

func (self *_Memory) ReadFloat32Pack(pd PackedData) (float32, error) {
	return self.ReadFloat32(Second(UnpackUI64[ValueType](uint64(pd))))
}

func (self *_Memory) ReadFloat64(offset uint32) (float64, error) {
	data, ok := self.module.raw.Memory().ReadFloat64Le(offset)
	if !ok {
		err := fmt.Errorf("Memory.ReadFloat64(%d, %d) out of range of memory size %d", offset, 8, self.Size())
		self.logger.Error(err.Error())
		return 0, err
	}

	return data, nil
}

func (self *_Memory) ReadFloat64Pack(pd PackedData) (float64, error) {
	return self.ReadFloat64(Second(UnpackUI64[ValueType](uint64(pd))))
}

func (self *_Memory) ReadString(offset uint32, size uint32) (string, error) {
	buf, err := self.ReadBytes(offset, size)
	if err != nil {
		return "", err
	}

	return string(buf), err
}

func (self *_Memory) ReadStringPack(pd PackedData) (string, error) {
	_, offset, size := UnpackUI64[ValueType](uint64(pd))
	return self.ReadString(offset, size)
}

// WriteAny writes a value of type interface{} to the memory buffer managed by the wazeroMemory instance,
// starting at the given offset.
//
// The method identifies the type of the value and performs the appropriate write operation.
func (self *_Memory) WriteAny(offset uint32, v any) error {
	var err error

	switch vTyped := v.(type) {
	case []byte:
		err = self.WriteBytes(offset, vTyped)
	case byte:
		err = self.WriteByte(offset, vTyped)
	case uint32:
		err = self.WriteUint32(offset, vTyped)
	case uint64:
		err = self.WriteUint64(offset, vTyped)
	case float32:
		err = self.WriteFloat32(offset, vTyped)
	case float64:
		err = self.WriteFloat64(offset, vTyped)
	case string:
		err = self.WriteString(offset, vTyped)
	default:
		err := fmt.Errorf("unsupported write data type %s", reflect.TypeOf(v))
		self.logger.Error(err.Error())
		return err
	}

	return err
}

func (self *_Memory) WriteBytes(offset uint32, v []byte) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().Write(offset, v),
		nil,
		fmt.Errorf("Memory.WriteBytes(%d, %d) out of range of memory size %d", offset, len(v), self.Size()),
	)))
}

func (self *_Memory) WriteBytesPack(v []byte) PackedData {
	size := uint32(len(v))

	offset, err := self.Malloc(size)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	err = self.WriteBytes(offset, v)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeBytes, offset, size)
}

func (self *_Memory) WriteByte(offset uint32, v byte) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteByte(offset, v),
		nil,
		fmt.Errorf("Memory.WriteByte(%d, %d) out of range of memory size %d", offset, 1, self.Size()),
	)))
}

func (self *_Memory) WriteBytePack(v byte) PackedData {
	offset, err := self.Malloc(1)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	err = self.WriteByte(offset, v)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeByte, offset, 1)
}

func (self *_Memory) WriteUint32(offset uint32, v uint32) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteUint32Le(offset, v),
		nil,
		fmt.Errorf("Memory.WriteUint32(%d, %d) out of range of memory size %d", offset, 4, self.Size()),
	)))
}

func (self *_Memory) WriteUint32Pack(v uint32) PackedData {
	offset, err := self.Malloc(4)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	err = self.WriteUint32(offset, v)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeI32, offset, 4)
}

func (self *_Memory) WriteUint64(offset uint32, v uint64) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteUint64Le(offset, v),
		nil,
		fmt.Errorf("Memory.WriteUint64(%d, %d) out of range of memory size %d", offset, 8, self.Size()),
	)))
}

func (self *_Memory) WriteUint64Pack(v uint64) PackedData {
	offset, err := self.Malloc(8)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	err = self.WriteUint64(offset, v)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeI64, offset, 8)
}

func (self *_Memory) WriteFloat32(offset uint32, v float32) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteFloat32Le(offset, v),
		nil,
		fmt.Errorf("Memory.WriteFloat32(%d, %d) out of range of memory size %d", offset, 8, self.Size()),
	)))
}

func (self *_Memory) WriteFloat32Pack(v float32) PackedData {
	offset, err := self.Malloc(4)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	err = self.WriteFloat32(offset, v)
	if err != nil {
		self.logger.Error(err.Error())
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeF32, offset, 4)
}

func (self *_Memory) WriteFloat64(offset uint32, v float64) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteFloat64Le(offset, v),
		nil,
		fmt.Errorf("Memory.WriteFloat64(%d, %d) out of range of memory size %d", offset, 8, self.Size()),
	)))
}

func (self *_Memory) WriteFloat64Pack(v float64) PackedData {
	offset, err := self.Malloc(8)
	if err != nil {
		return 0
	}

	err = self.WriteFloat64(offset, v)
	if err != nil {
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeF64, offset, 8)
}

func (self *_Memory) WriteString(offset uint32, v string) error {
	return Log(self.logger, Aggregate("failed to write data", Ternary(
		self.module.raw.Memory().WriteString(offset, v),
		nil,
		fmt.Errorf("Memory.WriteString(%d, %d) out of range of memory size %d", offset, len(v), self.Size()),
	)))
}

func (self *_Memory) WriteStringPack(v string) PackedData {
	size := uint32(len(v))

	offset, err := self.Malloc(size)
	if err != nil {
		return 0
	}

	err = self.WriteString(offset, v)
	if err != nil {
		return 0
	}

	return pack[PackedData](self.logger, types.ValueTypeString, offset, size)
}

func (self *_Memory) WriteMultiPack(datas ...PackedData) MultiPackedData {
	size := uint32(len(datas)) * 8
	if size == 0 {
		return 0
	}

	offset, err := self.Malloc(size)
	if err != nil {
		return 0
	}

	pdsU64 := Map(datas, func(data PackedData) uint64 { return uint64(data) })
	err = self.WriteBytes(offset, Uint64ArrayToBytes(pdsU64))
	if err != nil {
		return 0
	}

	return pack[MultiPackedData](self.logger, types.ValueTypeString, offset, size)
}

// Size returns the size in bytes available. e.g. If the underlying memory
// has 1 page: 65536
func (self *_Memory) Size() uint32 {
	return self.module.raw.Memory().Size()
}

// Malloc allocates memory in wasm linear memory with the specified size.
//
// It invokes the "malloc" GuestFunction of the associated wazeroModule using the provided size parameter.
// Returns the allocated memory offset and any encountered error.
//
// Malloc allows memory allocation from within a host function or externally,
// returning the allocated memory offset to be used in a guest function.
// This can be helpful, for instance, when passing string data from the host to the guest.
//
// NOTE: Always make sure to free memory after allocation.
func (self *_Memory) Malloc(size uint32) (uint32, error) {
	offset, err := self.module.GuestFunction(self.module.config.Context, "malloc").(*_GuestFunction).call(uint64(size))
	if err != nil {
		return 0, Log(self.logger, Aggregate("can't invoke 'malloc' function", err))
	}

	return uint32(offset), nil
}

// Free releases the memory block at the specified offset in wazeroMemory.
// It invokes the "free" GuestFunction of the associated wazeroModule using the provided offset parameter.
// Returns any encountered error during the memory deallocation.
func (self *_Memory) Free(offsets ...uint32) error {
	return Log(
		self.logger,
		Aggregate(
			"failed to invoke 'free' function",
			Map(offsets, func(offset uint32) error {
				return Second(self.module.GuestFunction(self.module.config.Context, "free").(*_GuestFunction).call(uint64(offset)))
			})...,
		),
	)
}

func (self *_Memory) FreePack(datas ...PackedData) error {
	return Log(
		self.logger,
		Aggregate(
			"failed to free up packed data",
			Map(datas, func(data PackedData) error {
				return self.Free(Second(UnpackUI64[ValueType](uint64(data))))
			})...,
		),
	)
}

func pack[T PackedData | MultiPackedData](logger logging.Logger, typ types.ValueType, offset uint32, size uint32) T {
	pd, err := PackUI64(typ, offset, size)
	if err != nil {
		_ = Log(logger, Aggregate("failed to pack data", err))
		return 0
	}

	return T(pd)
}
