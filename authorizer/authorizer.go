package authorizer

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
)

type Authorizer struct {
	env wasm.WasmEnv
	ptr uint64
}

func (authorizer Authorizer) New(env wasm.WasmEnv, ptr uint64) Authorizer {
	return Authorizer{env: env, ptr: ptr}
}

func (authorizer Authorizer) Ptr() uint64 {
	return authorizer.ptr
}

func (authorizer Authorizer) ToStringWasmFunction() string {
	return "authorizer_toString"
}

func (authorizer Authorizer) ToString() (string, error) {
	return shared.AsString(authorizer.env, authorizer)
}

func (authorizer Authorizer) Authorize() (uint64, error) {
	function, err := authorizer.env.GetFunction("authorizer_authorize")
	if err != nil {
		return 0, err
	}

	returnPtr, err := authorizer.env.GetReturnArea()
	defer authorizer.env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = authorizer.env.Call(function, returnPtr, authorizer.ptr)
	if err != nil {
		return 0, err
	}

	// Return the index of the matched policy
	matchedPolicy, err := authorizer.env.GetPointee(returnPtr)
	if err != nil {
		return 0, err
	}

	return matchedPolicy, nil
}
