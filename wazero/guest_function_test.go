package wazero_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	test_utils "github.com/wasify-io/wasify-go/internal/test-utils"
	"github.com/wasify-io/wasify-go/logging"
	"github.com/wasify-io/wasify-go/models"
)

func TestGuestFunctions(t *testing.T) {
	t.Run("successful instantiation", func(t *testing.T) {
		ctx := context.Background()

		runtime := test_utils.CreateRuntime(t, ctx, &models.RuntimeConfig{
			Runtime: models.RuntimeWazero,
			Logger:  logging.NewSlogLogger(logging.LogError),
		})

		module := test_utils.CreateModule(t, runtime, &models.ModuleConfig{
			Context:   context.Background(),
			Namespace: "guest_all_available_types",
			Logger:    logging.NewSlogLogger(logging.LogInfo),
			Wasm:      test_utils.LoadTestWASM(t, "guest_all_available_types"),
		})

		result := module.GuestFunction(ctx, "guestTest").Invoke(
			[]byte("bytes!"),
			byte(1),
			uint32(32),
			uint64(64),
			float32(32.0),
			float64(64.01),
			"Wasify",
			"any type",
		)
		assert.NoError(t, result.Error())

		t.Log("TestGuestFunctions RES:", result)
	})
}
