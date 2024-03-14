package models

import "io"

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
