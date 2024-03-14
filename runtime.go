package wasify

import (
	"context"
	"errors"

	"github.com/wasify-io/wasify-go/models"
	"github.com/wasify-io/wasify-go/wazero"
)

// NewRuntime creates and initializes a new runtime based on the provided configuration.
// It returns the initialized runtime and any error that might occur during the process.
func NewRuntime(ctx context.Context, config *models.RuntimeConfig) (models.Runtime, error) {
	config.Logger.Info("runtime has been initialized successfully", "runtime", config.Runtime)

	// Retrieve the appropriate runtime implementation based on the configured type.
	runtime := getRuntime(ctx, config)
	if runtime == nil {
		err := errors.New("unsupported runtime")
		config.Logger.Error(err.Error(), "runtime", config.Runtime)
		return nil, err
	}

	return runtime, nil
}

// getRuntime returns an instance of the appropriate runtime implementation
// based on the configured runtime type in the RuntimeConfig.
func getRuntime(ctx context.Context, config *models.RuntimeConfig) models.Runtime {
	switch config.Runtime {
	case models.RuntimeWazero:
		return wazero.NewRuntime(ctx, config)
	default:
		return nil
	}
}
