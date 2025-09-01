package biscuit

import "biscuit-wasm-go/wasm"

type Biscuit struct {
	env wasm.WasmEnv
	ptr uint64
}

func (biscuit Biscuit) New(env wasm.WasmEnv, ptr uint64) Biscuit {
	return Biscuit{env: env, ptr: ptr}
}
