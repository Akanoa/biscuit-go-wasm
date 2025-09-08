package authorizer

import (
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/wasm"
	"fmt"
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
	return "authorizer_print_world"
}

func (authorizer Authorizer) ToString() (string, error) {
	return shared.AsString(authorizer.env, authorizer)
}

func (authorizer Authorizer) Authorize() (uint32, error) {

	if authorizer.ptr == 0 {
		return 0, fmt.Errorf("authorizer not initialized")
	}

	returnArea, err := authorizer.env.GetReturnArea()
	if err != nil {
		return 0, err
	}

	_, err = authorizer.env.Call("authorizer_run_limits", returnArea, authorizer.ptr)
	if err != nil {
		return 0, err
	}

	runLimits, err := authorizer.env.ResultPointer(returnArea)
	if err != nil {
		return 0, err
	}

	_, err = authorizer.env.Call("authorizer_authorize", returnArea, authorizer.ptr, runLimits)
	if err != nil {
		return 0, err
	}

	// Return the index of the matched policy
	matchedPolicy, err := authorizer.env.ResultNumber(returnArea)
	fmt.Println("matched policy", matchedPolicy)
	fmt.Println("err", err)
	if err != nil {
		return 0, err
	}

	return matchedPolicy, nil
}
