package wazero

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	. "github.com/wasify-io/wasify-go/internal/utils"
	"github.com/wasify-io/wasify-go/logging"
	. "github.com/wasify-io/wasify-go/models"
)

// NewRuntime creates and returns a wazero runtime instance using the provided context and
// RuntimeConfig. It configures the runtime with specific settings and features.
func NewRuntime(ctx context.Context, config *RuntimeConfig) Runtime {
	// TODO: Allow user to control the following options:
	// 1. WithCloseOnContextDone
	// 2. Memory
	// Create a new wazero runtime instance with specified configuration options.
	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithCoreFeatures(api.CoreFeaturesV2).
		WithCustomSections(false).
		WithCloseOnContextDone(false).
		// Enable runtime debug if user sets LogSeverity to debug level in runtime configuration
		WithDebugInfoEnabled(config.Logger.Severity() == logging.LogDebug),
	)

	// Instantiate the runtime with the WASI snapshot preview1.
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return &_Runtime{
		config:  config,
		runtime: runtime,
	}
}

// The wazeroRuntime struct combines a wazero runtime instance with runtime configuration.
type _Runtime struct {
	runtime wazero.Runtime
	config  *RuntimeConfig
}

// Close closes the resource.
// Note: The context parameter is used for value lookup, such as for
// logging. A canceled or otherwise done context will not prevent Close
// from succeeding.
func (self *_Runtime) Close(ctx context.Context) error {
	return Log(
		self.config.Logger,
		Aggregate("failed to close runtime", self.runtime.Close(ctx)),
	)
}

// Create creates a new guest instance based on the provided ModuleConfig within
// the wazero runtime context. It returns the created guest and any potential error.
func (self *_Runtime) Create(config *ModuleConfig) (Module, error) {
	if err := self.verify(config); err != nil {
		return nil, Aggregate("failed to create module", err)
	}

	return _NewModule(self.runtime, config)
}

func (self *_Runtime) verify(config *ModuleConfig) error {
	// TODO - make sure that all required props are set in the config
	if config.Wasm.Hash == "" {
		return nil
	}

	actual, err := CalculateHash(config.Wasm.Binary)
	if err != nil {
		return Aggregate("failed to calculate hash for module", err)
	}

	return Aggregate("failed to check hash for module", CompareHashes(actual, config.Wasm.Hash))
}
