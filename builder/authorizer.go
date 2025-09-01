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
	function, err := env.GetFunction("authorizerbuilder_new")
	if err != nil {
		return AuthorizerBuilder{}, err
	}

	ret, err := env.Call(function)
	if err != nil {
		slog.Error("authorizerbuilder_new failed", slog.Any("err", err))
		return AuthorizerBuilder{}, err
	}

	builder.env = env
	builder.ptr = ret[0]

	return builder, nil
}

// Build builds the Biscuit using the provided private key.
func (builder *AuthorizerBuilder) Build(privateKey token.Biscuit) (authorizer.Authorizer, error) {
	if builder.ptr == 0 {
		return authorizer.Authorizer{}, fmt.Errorf("builder not initialized")
	}

	function, err := builder.env.GetFunction("authorizerbuilder_build")
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	returnPtr, err := builder.env.GetReturnArea()
	defer builder.env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = builder.env.Call(function, returnPtr, builder.ptr, privateKey.Ptr())
	if err != nil {
		slog.Error("authorizerbuilder_build failed", slog.Any("err", err))
		return authorizer.Authorizer{}, err
	}

	valuePtr, err := builder.env.GetPointee(returnPtr)
	if err != nil {
		slog.Error("authorizerbuilder_build failed, unable to get return value", slog.Any("err", err))
		return authorizer.Authorizer{}, err
	}

	return authorizer.Authorizer{}.New(builder.env, valuePtr), nil
}

// AddCode adds the provided code to the AuthorizerBuilder.
func (builder *AuthorizerBuilder) AddCode(code string) error {
	if builder.ptr == 0 {
		return fmt.Errorf("builder builder not initialized")
	}

	function, err := builder.env.GetFunction("authorizerbuilder_addCode")
	if err != nil {
		return err
	}

	strPtr, err := builder.env.WriteString(code)
	if err != nil {
		return fmt.Errorf("cannot write string to wasm memory: %w", err)
	}
	defer builder.env.Free(strPtr, uint64(len(code)))

	returnPtr, err := builder.env.GetReturnArea()
	if err != nil {
		return err
	}

	_, err = builder.env.Call(function, returnPtr, builder.ptr, strPtr, uint64(len(code)))
	if err != nil {
		return fmt.Errorf("authorizerbuilder_addCode failed: %w", err)
	}
	return nil
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (builder AuthorizerBuilder) ToString() (string, error) {
	if builder.ptr == 0 {
		return "", fmt.Errorf("token builder not initialized")
	}

	return shared.AsString(builder.env, builder)
}
