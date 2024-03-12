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
	ReadPacks() ([]PackedData, error)
}

func NewGuestFunctionResult(data uint64, memory Memory) GuestFunctionResult {
	return &_GuestFunctionResult{
		data:   data,
		memory: memory,
	}
}

type _GuestFunctionResult struct {
	data   uint64
	memory Memory
	result []PackedData
}

func (self *_GuestFunctionResult) Close() error {
	return self.memory.FreePack(self.result...)
}

// ReadPacks decodes the packedData from a GuestFunctionResult instance and retrieves a sequence of packed datas.
// NOTE: Frees multiPackedData, which means ReadPacks should be called once.
func (self *_GuestFunctionResult) ReadPacks() ([]PackedData, error) {

	if self.data == 0 {
		return nil, errors.New("packedData is empty")
	}

	t, offsetU32, size := utils.UnpackUI64(self.data)

	if t != types.ValueTypePack {
		err := fmt.Errorf("Can't unpack host data, the type is not a valueTypePack. expected %d, got %d", types.ValueTypePack, t)
		return nil, err
	}

	bytes, err := self.memory.ReadBytes(offsetU32, size)
	if err != nil {
		err := errors.Join(errors.New("ReadPacks error, can't read bytes:"), err)
		return nil, err
	}

	err = self.memory.FreePack(PackedData(self.data))
	if err != nil {
		err := errors.Join(errors.New("ReadPacks error, can't free multiPackedData:"), err)
		return nil, err
	}

	self.result = utils.Map(
		utils.BytesToUint64Array(bytes),
		func(data uint64) PackedData {
			return PackedData(data)
		},
	)

	return self.result, nil
}
