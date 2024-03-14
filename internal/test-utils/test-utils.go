package test_utils

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wasify-io/wasify-go"
	"github.com/wasify-io/wasify-go/models"
)

func CreateRuntime(t *testing.T, ctx context.Context, config *models.RuntimeConfig) models.Runtime {
	runtime, err := wasify.NewRuntime(ctx, config)

	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, runtime.Close(ctx))
	})

	return runtime
}

func CreateModule(t *testing.T, runtime models.Runtime, config *models.ModuleConfig) models.Module {
	module, err := runtime.Create(config)

	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, module.Close(config.Context))
	})

	return module
}

func LoadTestWASM(t *testing.T, path string) models.Wasm {
	data, err := os.ReadFile(fmt.Sprintf("../resources/testdata/wasm/%s/main.wasm", path))

	assert.NoError(t, err)

	return models.Wasm{
		Binary: data,
	}
}
