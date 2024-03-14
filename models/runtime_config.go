package models

import "github.com/wasify-io/wasify-go/logging"

// The RuntimeConfig struct holds configuration settings for a runtime.
type RuntimeConfig struct {
	// Specifies the type of runtime being used.
	Runtime RuntimeType
	// Logger to use for the runtime and module
	Logger logging.Logger
}
