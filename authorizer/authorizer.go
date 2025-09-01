package authorizer

import "biscuit-wasm-go/wasm"

type Authorizer struct {
	env wasm.WasmEnv
	ptr uint64
}

func (authorizer Authorizer) New(env wasm.WasmEnv, ptr uint64) Authorizer {
	return Authorizer{env: env, ptr: ptr}
}
