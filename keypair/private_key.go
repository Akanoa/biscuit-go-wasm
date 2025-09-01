package keypair

import (
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

type PrivateKey struct {
	env wasm.WasmEnv
	ptr uint64
}

func (private_key PrivateKey) Ptr() uint64 {
	return private_key.ptr
}

func (private_key PrivateKey) New(env wasm.WasmEnv) PrivateKey {
	return PrivateKey{env: env, ptr: 0}
}

func (self PrivateKey) ToString() (string, error) {
	if self.ptr == 0 {
		slog.Error("private key not initialized")
		return "", fmt.Errorf("private key not initialized")
	}

	function, err := self.env.GetFunction("privatekey_toString")
	if err != nil {
		slog.Error("exported function 'privatekey_toString' not found")
		return "", err
	}

	resultPtr, err := self.env.GetStringArea()
	if err != nil {
		slog.Error("malloc failed", slog.Any("err", err))
		return "", err
	}
	defer self.env.Free(resultPtr, wasm.StringAreaSize)

	_, err = self.env.Call(function, resultPtr, self.ptr)
	if err != nil {
		slog.Error("privatekey_toString failed", slog.Any("err", err))
		return "", err
	}

	return self.env.GetStringValueFromPointer(resultPtr)
}

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
