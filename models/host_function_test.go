package models_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	test_utils "github.com/wasify-io/wasify-go/internal/test-utils"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

func TestHostFunctions(t *testing.T) {
	t.Run("successful instantiation", func(t *testing.T) {
		ctx := context.Background()

		runtime := test_utils.CreateRuntime(t, ctx, &RuntimeConfig{
			Runtime: RuntimeWazero,
			Logger:  logging.NewSlogLogger(logging.LogInfo),
		})

		module := test_utils.CreateModule(t, runtime, &ModuleConfig{
			Context:   context.Background(),
			Logger:    logging.NewSlogLogger(logging.LogInfo),
			Namespace: "host_all_available_types",
			Wasm:      test_utils.LoadTestWASM(t, "host_all_available_types"),
			HostFunctions: []HostFunction{
				{
					Name: "hostTest",
					Callback: func(ctx context.Context, module Module, params []PackedData) MultiPackedData {
						memory := module.Memory()

						_bytes, _ := memory.ReadBytesPack(params[0])
						assert.Equal(t, []byte("Guest: Wello Wasify!"), _bytes)

						_byte, _ := memory.ReadBytePack(params[1])
						assert.Equal(t, byte(1), _byte)

						_uint32, _ := memory.ReadUint32Pack(params[2])
						assert.Equal(t, uint32(11), _uint32)

						_uint64, _ := memory.ReadUint64Pack(params[3])
						assert.Equal(t, uint64(2023), _uint64)

						_float32, _ := memory.ReadFloat32Pack(params[4])
						assert.Equal(t, float32(11.1), _float32)

						_float64, _ := memory.ReadFloat64Pack(params[5])
						assert.Equal(t, float64(11.2023), _float64)

						_string, _ := memory.ReadStringPack(params[6])
						assert.Equal(t, "Guest: Wasify.", _string)

						return memory.WriteMultiPack(
							memory.WriteBytesPack([]byte("Some")),
							memory.WriteBytePack(1),
							memory.WriteUint32Pack(11),
							memory.WriteUint64Pack(2023),
							memory.WriteFloat32Pack(11.1),
							memory.WriteFloat64Pack(11.2023),
							memory.WriteStringPack("Host: Wasify."),
						)

					},
					Params: []ValueType{
						ValueTypeBytes,
						ValueTypeByte,
						ValueTypeI32,
						ValueTypeI64,
						ValueTypeF32,
						ValueTypeF64,
						ValueTypeString,
					},
					Results: []ValueType{
						ValueTypeBytes,
						ValueTypeByte,
						ValueTypeI32,
						ValueTypeI64,
						ValueTypeF32,
						ValueTypeF64,
						ValueTypeString,
					},
				},
			},
		})

		result := module.GuestFunction(ctx, "guestTest").Invoke()

		assert.NoError(t, result.Error())

		t.Log("TestHostFunctions RES:", result)
	})
}
