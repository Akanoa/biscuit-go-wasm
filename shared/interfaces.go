// Package shared contains cross-package helpers and small interfaces.
package shared

// Stringable abstracts types that can be turned into a string via a specific
// wasm export and that expose a raw pointer to their underlying WASM object.
type Stringable interface {
	ToStringWasmFunction() string
	Ptr() uint64
}
