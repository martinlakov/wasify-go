package wasify

import (
	"context"

	"github.com/wasify-io/wasify-go/internal/utils"
)

const WASIFY_NAMESPACE = "wasify"

// hostFunctions is a list of pre-defined host functions
type hostFunctions struct {
	config *ModuleConfig
}

func newHostFunctions(config *ModuleConfig) *hostFunctions {
	return &hostFunctions{config}
}

// newLog logs data from the guest module to the host machine,
// to avoid stdin/stdout calls and ensure sandboxing.
func (hf *hostFunctions) newLog() *HostFunction {
	return &HostFunction{
		moduleConfig: hf.config,

		Name:    "log",
		Params:  []ValueType{ValueTypeString, ValueTypeBytes},
		Results: nil,
		Callback: func(ctx context.Context, m *ModuleProxy, params []PackedData) MultiPackedData {
			hf.config.log.Log(
				utils.LogSeverity(utils.Must(m.Memory.ReadBytePack(params[1]))),
				utils.Must(m.Memory.ReadStringPack(params[0])),
			)

			return 0
		},
	}
}
