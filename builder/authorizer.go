package builder

import (
	"biscuit-wasm-go/authorizer"
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

type AuthorizerBuilder struct {
	env wasm.WasmEnv
	ptr uint64
}

func (builder AuthorizerBuilder) Ptr() uint64 {
	return builder.ptr
}

func (builder AuthorizerBuilder) ToStringWasmFunction() string {
	return "authorizerbuilder_toString"
}

// New creates a new AuthorizerBuilder using a Wasm environment.
// Returns the initialized AuthorizerBuilder or an error in case of failure.
func (builder AuthorizerBuilder) New(env wasm.WasmEnv) (AuthorizerBuilder, error) {

	returnArea, err := env.GetReturnArea()
	if err != nil {
		return AuthorizerBuilder{}, err
	}

	_, err = env.Call("authorizer_builder_new", returnArea)
	if err != nil {
		slog.Error("authorizerbuilder_new failed", slog.Any("err", err))
		return AuthorizerBuilder{}, err
	}

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return AuthorizerBuilder{}, err
	}

	builder.ptr = ptr
	builder.env = env

	return builder, nil
}

// Build builds the Biscuit using the provided private key.
func (builder *AuthorizerBuilder) Build(biscuit token.Biscuit) (authorizer.Authorizer, error) {
	if builder.ptr == 0 {
		return authorizer.Authorizer{}, fmt.Errorf("builder not initialized")
	}

	returnArea, err := builder.env.GetReturnArea()
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	_, err = builder.env.Call("authorizer_builder_build", returnArea, builder.ptr, biscuit.Ptr())
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	ptr, err := builder.env.ResultPointer(returnArea)
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	return authorizer.Authorizer{}.New(builder.env, ptr), nil
}

// AddCode adds the provided code to the AuthorizerBuilder.
func (builder *AuthorizerBuilder) AddCode(code string) error {
	if builder.ptr == 0 {
		return fmt.Errorf("builder builder not initialized")
	}

	returnArea, err := builder.env.GetReturnArea()
	if err != nil {
		return err
	}

	strPtr, err := builder.env.WriteBytesToWasm([]byte(code))
	if err != nil {
		return err
	}
	defer builder.env.Free(strPtr, uint64(len(code)))

	_, err = builder.env.Call("authorizer_builder_add_code", returnArea, builder.ptr, strPtr, uint64(len(code)))
	if err != nil {
		return fmt.Errorf("authorizerbuilder_addCode failed: %w", err)
	}

	_, err = builder.env.ResultPointer(returnArea)
	return err
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (builder AuthorizerBuilder) ToString() (string, error) {
	if builder.ptr == 0 {
		return "", fmt.Errorf("token builder not initialized")
	}

	return shared.AsString(builder.env, builder)
}
