package models

type RMemory interface {
	ReadBytes(offset uint32, size uint32) ([]byte, error)
	ReadByte(offset uint32) (byte, error)
	ReadUint32(offset uint32) (uint32, error)
	ReadUint64(offset uint32) (uint64, error)
	ReadFloat32(offset uint32) (float32, error)
	ReadFloat64(offset uint32) (float64, error)
	ReadString(offset uint32, size uint32) (string, error)

	ReadAnyPack(pd PackedData) (any, uint32, uint32, error)
	ReadBytesPack(pd PackedData) ([]byte, error)
	ReadBytePack(pd PackedData) (byte, error)
	ReadUint32Pack(pd PackedData) (uint32, error)
	ReadUint64Pack(pd PackedData) (uint64, error)
	ReadFloat32Pack(pd PackedData) (float32, error)
	ReadFloat64Pack(pd PackedData) (float64, error)
	ReadStringPack(pd PackedData) (string, error)

	Size() uint32
	Free(...uint32) error
	FreePack(...PackedData) error
}

type WMemory interface {
	WriteAny(offset uint32, v any) error
	WriteBytes(offset uint32, v []byte) error
	WriteByte(offset uint32, v byte) error
	WriteUint32(offset uint32, v uint32) error
	WriteUint64(offset uint32, v uint64) error
	WriteFloat32(offset uint32, v float32) error
	WriteFloat64(offset uint32, v float64) error
	WriteString(offset uint32, v string) error

	WriteBytesPack(v []byte) PackedData
	WriteBytePack(v byte) PackedData
	WriteUint32Pack(v uint32) PackedData
	WriteUint64Pack(v uint64) PackedData
	WriteFloat32Pack(v float32) PackedData
	WriteFloat64Pack(v float64) PackedData
	WriteStringPack(v string) PackedData

	WriteMultiPack(...PackedData) MultiPackedData

	Malloc(size uint32) (uint32, error)
}

type Memory interface {
	RMemory
	WMemory
}
