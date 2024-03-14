package models

import (
	"context"

	"github.com/wasify-io/wasify-go/logging"
)

type ModuleConfig struct {
	// Module Namespace. Required.
	Namespace string

	// FSConfig configures a directory to be pre-opened for access by the WASI module if Enabled is set to true.
	// If GuestDir is not provided, the default guest directory will be "/".
	// Note: If FSConfig is not provided or Enabled is false, the directory will not be attached to WASI.
	FSConfig FSConfig

	// WASM configuration. Required.
	Wasm Wasm

	// List of host functions to be registered.
	HostFunctions []HostFunction

	Context context.Context

	Logger logging.Logger
}
