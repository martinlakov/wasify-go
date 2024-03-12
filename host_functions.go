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
		Params:  []ValueType{ValueTypeBytes, ValueTypeString},
		Results: nil,
		Callback: func(ctx context.Context, m *ModuleProxy, params []PackedData) MultiPackedData {
			severity := utils.Must(m.Memory.ReadBytePack(params[0]))
			level := utils.GetlogLevel(utils.LogSeverity(severity))
			message := utils.Must(m.Memory.ReadStringPack(params[1]))

			hf.config.log.Log(ctx, level, message)

			return 0
		},
	}
}
