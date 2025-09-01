package keypair

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
)

// PrivateKey represents a WASM private key.
type PrivateKey struct {
	env wasm.WasmEnv
	ptr uint64
}

// ToStringWasmFunction returns the name of the WASM function that converts the PublicKey to its string representation.
func (private_key PrivateKey) ToStringWasmFunction() string {
	return "privatekey_toString"
}

// Ptr returns the pointer to the underlying WASM object.
func (private_key PrivateKey) Ptr() uint64 {
	return private_key.ptr
}

// New creates a new PrivateKey using the linked Wasm environment.
func (private_key PrivateKey) New(env wasm.WasmEnv) PrivateKey {
	return PrivateKey{env: env, ptr: 0}
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (private_key PrivateKey) ToString() (string, error) {
	if private_key.ptr == 0 {
		return "", fmt.Errorf("biscuit private_key not initialized")
	}

	return shared.AsString(private_key.env, private_key)
}

// FromString converts the string representation of a PrivateKey to a PrivateKey using the linked Wasm environment. Returns an error.
func (self *PrivateKey) FromString(data string) error {
	// Note: Go strings are UTF-8 already. We must copy bytes into WASM memory
	// and pass (ptr, len) according to wasm-bindgen ABI.
	function, err := self.env.GetFunction("privatekey_fromString")
	if err != nil {
		return err
	}

	strPtr, err := self.env.WriteString(data)
	if err != nil {
		return fmt.Errorf("cannot write string to wasm memory: %w", err)
	}
	defer self.env.Free(strPtr, uint64(len(data)))

	// Allocate return area (3 u32 values: value_ptr, error_ptr, is_err)
	retPtr, err := self.env.GetReturnArea()
	if err != nil {
		return fmt.Errorf("malloc for return area failed: %w", err)
	}
	defer self.env.Free(retPtr, wasm.ReturnAreaSize)

	// Call: privatekey_fromString(out_ptr, str_ptr, str_len)
	_, err = self.env.Call(function, retPtr, strPtr, uint64(len(data)))
	if err != nil {
		return fmt.Errorf("privatekey_fromString failed: %w", err)
	}

	// Read result triple
	valuePtr, err := self.env.GetPointee(retPtr)
	if err != nil {
		return err
	}

	self.ptr = valuePtr
	return nil
}
