package wazero

import (
	"context"
	"errors"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	. "github.com/wasify-io/wasify-go/internal/utils"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

func _NewModule(runtime wazero.Runtime, config *ModuleConfig) (*_Module, error) {
	result := &_Module{config: config}

	return result, result.initialize(runtime)
}

// The wazeroModule struct combines an instantiated wazero modul
// with the generic guest configuration.
type _Module struct {
	raw    api.Module
	config *ModuleConfig
	closer func(context.Context) error
}

// Memory retrieves a Memory instance associated with the wazeroModule.
func (self *_Module) Memory() Memory {
	return &_Memory{
		module: self,
		logger: self.config.Logger,
	}
}

// Close closes the resource.
//
// Note: The context parameter is used for value lookup, such as for
// logging. A canceled or otherwise done context will not prevent Close
// from succeeding.
func (self *_Module) Close(ctx context.Context) error {
	// TODO - invoke the closer too
	err := self.raw.Close(ctx)
	if err != nil {
		err = errors.Join(errors.New("can't close guest"), err)
		self.config.Logger.Error(err.Error())
		return err
	}

	return nil
}

// GuestFunction returns a GuestFunction instance associated with the wazeroModule.
// GuestFunction is used to work with exported function from this guest.
//
// Example usage:
//
//	result = guest.GuestFunction(ctx, "greet").Invoke("argument1", "argument2", 123)
//	if err := result.Error(); err != nil {
//	    slog.Error(err.Error())
//	}
func (self *_Module) GuestFunction(ctx context.Context, name string) GuestFunction {
	fn := self.raw.ExportedFunction(name)
	if fn == nil {
		self.config.Logger.Warn("exported function does not exist", "function", name, "namespace", self.config.Namespace)
	}

	return &_GuestFunction{
		name:   name,
		ctx:    ctx,
		fn:     fn,
		config: self.config,
		memory: self.Memory(),
	}
}

func (self *_Module) initialize(runtime wazero.Runtime) error {
	defaults, err := self.instantiateDefaultHostFunctions(runtime)
	if err != nil {
		return Aggregate("failed to instantiate default host functions", err)
	}

	custom, err := self.instantiateUserDefinedHostFunctions(runtime)
	if err != nil {
		return Aggregate("failed to instantiate custom host functions", err)
	}

	guest, err := self.instantiateGuestFunctions(runtime)
	if err != nil {
		return Aggregate("failed to instantiate guest functions", err)
	}

	self.raw = guest
	self.closer = func(ctx context.Context) error {
		return Aggregate("failed to close module", guest.Close(ctx), custom.Close(ctx), defaults.Close(ctx))
	}

	return nil
}

// instantiateModule compiles and instantiates a WebAssembly guest using the wazero runtime.
// It compiles the guest, creates a guest configuration, and then instantiates the guest.
// Returns the instantiated guest and any potential error.
func (self *_Module) instantiateGuestFunctions(runtime wazero.Runtime) (api.Module, error) {
	// TODO: Add more configurations
	cfg := wazero.NewModuleConfig()
	cfg = cfg.WithStdin(os.Stdin)
	cfg = cfg.WithStdout(os.Stdout)
	cfg = cfg.WithStderr(os.Stderr)

	if self.config.FSConfig.Enabled {
		cfg = cfg.WithFSConfig(
			wazero.NewFSConfig().
				WithDirMount(self.config.FSConfig.HostDir, self.config.FSConfig.GetGuestDir()),
		)
	}

	// Instantiate the compiled guest with the provided guest configuration.
	module, err := runtime.InstantiateWithConfig(self.config.Context, self.config.Wasm.Binary, cfg)

	return module, Aggregate("failed to instantiate module guest functions", err)
}

// instantiateHostFunctions sets up and exports host functions for the guest using the wazero runtime.
// It configures host function callbacks, data types, and exports.
func (self *_Module) instantiateDefaultHostFunctions(runtime wazero.Runtime) (api.Module, error) {
	builder := runtime.NewHostModuleBuilder("wasify")

	for _, function := range self.predefined() {
		builder = builder.
			NewFunctionBuilder().
			WithGoModuleFunction(
				api.GoModuleFunc(wazeroHostFunctionCallback(self.config.Logger, self, function)),
				self.convertToAPIValueTypes(function.Params),
				self.convertToAPIValueTypes(function.Results),
			).
			Export(function.Name)
	}

	result, err := builder.Instantiate(self.config.Context)

	return result, Aggregate("failed to instantiate predefined host functions", err)
}

// instantiateHostFunctions sets up and exports host functions for the guest using the wazero runtime.
// It configures host function callbacks, data types, and exports.
func (self *_Module) instantiateUserDefinedHostFunctions(runtime wazero.Runtime) (api.Module, error) {
	config := self.config
	logger := self.config.Logger
	builder := runtime.NewHostModuleBuilder(config.Namespace)

	for _, function := range config.HostFunctions {
		logger.Debug("build host function", "namespace", config.Namespace, "function", function.Name)

		// Associate the host function with guest-related information.
		// This configuration ensures that the host function can access ModuleConfig data from various contexts.
		// See host_function.go for more details.
		function.Config = config

		// If host function has any return values, we pack it as a single uint64
		var resultValuesPackedData = make([]ValueType, 0)
		if len(function.Results) > 0 {
			resultValuesPackedData = []ValueType{ValueTypeI64}
		}

		builder = builder.
			NewFunctionBuilder().
			WithGoModuleFunction(
				api.GoModuleFunc(wazeroHostFunctionCallback(logger, self, &function)),
				self.convertToAPIValueTypes(function.Params),
				self.convertToAPIValueTypes(resultValuesPackedData),
			).
			Export(function.Name)
	}

	result, err := builder.Instantiate(config.Context)

	return result, Aggregate("failed to instantiate user-defined host functions", err)
}

func (self *_Module) predefined() []*HostFunction {
	return []*HostFunction{
		{
			Config: self.config,

			Name:    "log",
			Params:  []ValueType{ValueTypeString, ValueTypeBytes},
			Results: nil,
			Callback: func(ctx context.Context, module Module, params []PackedData) MultiPackedData {
				memory := module.Memory()
				self.config.Logger.Log(
					logging.LogSeverity(Must(memory.ReadBytePack(params[1]))),
					Must(memory.ReadStringPack(params[0])),
				)

				return 0
			},
		},
	}
}

// convertToAPIValueTypes converts an array of ValueType values to their corresponding
// api.ValueType representations used by the Wazero runtime.
// ValueType describes a parameter or result type mapped to a WebAssembly
// function signature.
func (self *_Module) convertToAPIValueTypes(types []ValueType) []api.ValueType {
	valueTypes := make([]api.ValueType, len(types))

	for i, t := range types {
		switch t {
		case
			ValueTypeBytes,
			ValueTypeByte,
			ValueTypeI32,
			ValueTypeI64,
			ValueTypeF32,
			ValueTypeF64,
			ValueTypeString:
			valueTypes[i] = api.ValueTypeI64
		}
	}

	return valueTypes
}
