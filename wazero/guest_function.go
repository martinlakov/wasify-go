package wazero

import (
	"context"
	"errors"
	"fmt"

	"github.com/tetratelabs/wazero/api"
	"github.com/wasify-io/wasify-go/internal/types"
	. "github.com/wasify-io/wasify-go/internal/utils"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

type _GuestFunction struct {
	name   string
	ctx    context.Context
	fn     api.Function
	memory Memory
	config *ModuleConfig
}

// Invoke calls a specified guest function with the provided parameters. It ensures proper memory management,
// data conversion, and compatibility with data types. Each parameter is converted to its packedData format,
// which provides a compact representation of its memory offset, size, and type information. This packedData
// is written into the WebAssembly memory, allowing the guest function to correctly interpret and use the data.
//
// While the method takes care of memory allocation for the parameters and writing them to memory, it does
// not handle freeing the allocated memory. If an error occurs at any step, from data conversion to memory
// allocation, or during the guest function invocation, the error is logged, and the function returns with an error.
//
// Example:
//
// res := guest.GuestFunction(ctx, "guestTest").Invoke([]byte("bytes!"), uint32(32), float32(32.0), "Wasify")
//
// params ...any: A variadic list of parameters of any type that the user wants to pass to the guest function.
//
// Return value: The result of invoking the guest function in the form of a GuestFunctionResult pointer,
// or an error if any step in the process fails.
func (self *_GuestFunction) Invoke(params ...any) GuestFunctionResult {
	self.config.Logger.Log(
		self.severity(),
		"calling guest function", "namespace", self.config.Namespace, "function", self.name, "params", params,
	)

	stack, err := self.process(make([]uint64, len(params)), params...)
	if err != nil {
		return self.errorResult(Aggregate(fmt.Sprintf("failed to process parameters for guest function %s", self.name), err))
	}

	data, err := self.call(stack...)
	if err != nil {
		return self.errorResult(Aggregate(fmt.Sprintf("failed to invoke the guest function '%s'", self.name), err))
	}

	return _NewGuestFunctionResult(nil, data, self.memory)
}

func (self *_GuestFunction) severity() logging.LogSeverity {
	return Ternary(
		OneOf(self.config.Namespace, "malloc", "free"),
		logging.LogDebug,
		logging.LogInfo,
	)
}

func (self *_GuestFunction) errorResult(err error) GuestFunctionResult {
	return _NewGuestFunctionResult(Log(self.config.Logger, err), 0, nil)
}

// Call invokes wazero's CallWithStack method, which returns ome uint64 message,
// in most cases it is used to call built in methods such as "malloc", "free"
// See wazero's CallWithStack for more details.
func (self *_GuestFunction) call(params ...uint64) (MultiPackedData, error) {
	// size of params len(params) + one size for return uint64 value
	stack := make([]uint64, len(params)+1)
	copy(stack, params)

	err := self.fn.CallWithStack(self.ctx, stack[:])
	if err != nil {
		err = errors.Join(errors.New("error invoking internal call func"), err)
		self.config.Logger.Error(err.Error())
		return 0, err
	}

	return MultiPackedData(stack[0]), nil
}

func (self *_GuestFunction) process(stack []uint64, params ...any) ([]uint64, error) {
	for i, p := range params {
		valueType, offsetSize, err := types.GetOffsetSizeAndDataTypeByConversion(p)
		if err != nil {
			return stack, Aggregate(fmt.Sprintf("failed to convert guest function parameter %s", self.name), err)
		}

		// allocate memory for each value
		offsetI32, err := self.memory.Malloc(offsetSize)
		if err != nil {
			return stack, Aggregate(fmt.Sprintf("an error occurred while attempting to alloc memory for guest func param in: %s", self.name), err)
		}

		if err = self.memory.WriteAny(offsetI32, p); err != nil {
			return stack, Aggregate("failed to write arg to memory", err)
		}

		stack[i], err = PackUI64(valueType, offsetI32, offsetSize)
		if err != nil {
			return stack, Aggregate(fmt.Sprintf("failed to pack data for guest func param in:  %s", self.name), err)
		}
	}

	return stack, nil
}
