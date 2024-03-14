package wazero

import (
	"context"

	"github.com/tetratelabs/wazero/api"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

// wazeroHostFunctionCallback returns a callback function that acts as a bridge between
// the host function and the wazero runtime. This bridge ensures the seamless integration of
// the host function within the wazero environment by managing various tasks, including:
//
//   - Initialization of wazeroModule and ModuleProxy to set up the execution environment.
//   - Converting stack parameters into structured parameters that the host function can understand.
//   - Executing the user-defined host function callback with the correctly formatted parameters.
//   - Processing the results of the host function, converting them back into packed data format,
//     and writing the final packed data into linear memory.
//
// Diagram of the wazeroHostFunctionCallback Process:
// +--------------------------------------+
// | wazeroHostFunctionCallback           |
// |                                      |
// |  +---------------------------+       |
// |  | Initialize wazeroModule   |       |
// |  | and ModuleProxy           |       |
// |  +---------------------------+       |
// |                                      |
// |  +----------------------------+      |
// |  | Convert Stack Params to    |      |
// |  | Structured Params for      |      |
// |  | Host Function              |      |
// |  +----------------------------+      |
// |                 |                    |
// |                 v                    |
// |  +----------------------------+      |
// |  | 🚀 Execute User-defined    |      |
// |  | Host Function Callback     |      |
// |  +----------------------------+      |
// |                 |                    |
// |                 v                    |
// |  +-----------------------------+     |
// |  | Convert Return Values to    |     |
// |  | Packed Data using           |     |
// |  | writeResultsToMemory        |     |
// |  | and write final packedData  |     |
// |  | into linear memory          |     |
// |  +-----------------------------+     |
// |                                      |
// +--------------------------------------+
//
// Return value: A callback function that takes a context, api.Module, and a stack of parameters,
// and handles the integration of the host function within the wazero runtime.
func wazeroHostFunctionCallback(logger logging.Logger, module *_Module, function *HostFunction) func(context.Context, api.Module, []uint64) {
	return func(ctx context.Context, mod api.Module, stack []uint64) {
		params, err := function.PreHostFunctionCallback(stack)
		if err != nil {
			logger.Error(err.Error(), "namespace", module.config.Namespace, "func", function.Name)
		}

		results := function.Callback(ctx, module, params)

		function.PostHostFunctionCallback(results, stack)
	}
}
