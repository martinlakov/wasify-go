package models

// RuntimeType defines a type of WebAssembly (wasm) runtime.
//
// Currently, the only supported wasm runtime is Wazero.
// However, in the future, more runtimes could be added.
// This means that you'll be able to run modules
// on various wasm runtimes.
type RuntimeType uint8

const (
	RuntimeWazero RuntimeType = iota
)

func (rt RuntimeType) String() (runtimeName string) {

	switch rt {
	case RuntimeWazero:
		runtimeName = "Wazero"
	}

	return
}
