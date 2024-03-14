package models

// Wasm configures a new wasm file.
// Binay is required.
// Hash is optional.
type Wasm struct {
	Binary []byte
	Hash   string
}
