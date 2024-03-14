package models

import (
	"context"
	"fmt"

	"github.com/wasify-io/wasify-go/internal/utils"
)

// HostFunction defines a host function that can be invoked from a guest module.
type HostFunction struct {
	// Callback function to execute when the host function is invoked.
	Callback HostFunctionCallback

	// Name of the host function.
	Name string

	// Params specifies the types of parameters that the host function expects.
	//
	// The length of 'Params' should match the expected number of arguments
	// from the host function when called from the guest.
	Params []ValueType

	// Results specifies the types of values that the host function Results.
	//
	// The length of 'Results' should match the expected number of Results
	// from the host function as used in the guest.
	Results []ValueType

	// Allocation map to track parameter and return value allocations for host func.

	// Configuration of the associated module.
	Config *ModuleConfig
}

// HostFunctionCallback is the function signature for the callback executed by a host function.
//
// HostFunctionCallback encapsulates the runtime's internal implementation details.
// It serves as an intermediary invoked between the processing of function parameters and the final return of the function.
type HostFunctionCallback func(ctx context.Context, module Module, datas []PackedData) MultiPackedData

// PreHostFunctionCallback
// prepares parameters for the host function by converting
// packed stack parameters into a slice of PackedData. It validates parameter counts
// and leverages ModuleProxy for reading the data.
func (hf *HostFunction) PreHostFunctionCallback(stackParams []uint64) ([]PackedData, error) {
	// If user did not define params, skip the whole process, we still might get stackParams[0] = 0
	if len(hf.Params) == 0 {
		return nil, nil
	}

	if len(hf.Params) != len(stackParams) {
		return nil, fmt.Errorf("%s: params mismatch expected: %d received: %d ", hf.Name, len(hf.Params), len(stackParams))
	}

	return utils.Map(stackParams, func(param uint64) PackedData {
		return PackedData(param)
	}), nil
}

// PostHostFunctionCallback
// stores the resulting MultiPackedData into linear memory after the host function execution.
func (hf *HostFunction) PostHostFunctionCallback(mpd MultiPackedData, stackParams []uint64) {
	// Store final MultiPackedData into linear memory
	stackParams[0] = uint64(mpd)
}
