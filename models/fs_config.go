package models

import . "github.com/wasify-io/wasify-go/internal/utils"

// FSConfig configures a directory to be pre-opened for access by the WASI module if Enabled is set to true.
// If GuestDir is not provided, the default guest directory will be "/".
// Note: If FSConfig is not provided or Enabled is false, the directory will not be attached to WASI.
type FSConfig struct {
	// Whether to Enabled the directory for WASI access.
	Enabled bool

	// The directory on the host system.
	// Default: "/"
	HostDir string

	// The directory accessible to the WASI module.
	GuestDir string
}

// GetGuestDir gets the default path for guest module.
func (fs *FSConfig) GetGuestDir() string {
	return Ternary(fs.GuestDir == "", "/", fs.GuestDir)
}
