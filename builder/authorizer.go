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

	if len(ret) == 0 {
		return AuthorizerBuilder{}, fmt.Errorf("no result returned from authorizerbuilder_new")
	}

	builder.env = env
	builder.ptr = ret[0]

	fmt.Println("AuthorizerBuilder ptr:", builder.ptr)

	return builder, nil
}

// Build builds the Biscuit using the provided private key.
func (builder *AuthorizerBuilder) Build(biscuit token.Biscuit) (authorizer.Authorizer, error) {
	if builder.ptr == 0 {
		return authorizer.Authorizer{}, fmt.Errorf("builder not initialized")
	}

	// First try the authenticated build path.
	function, err := builder.env.GetFunction("authorizerbuilder_buildAuthenticated")
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	returnPtr, err := builder.env.GetReturnArea()
	if err != nil {
		return authorizer.Authorizer{}, err
	}
	defer builder.env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = builder.env.Call(function, returnPtr, builder.ptr, biscuit.Ptr())
	if err != nil {
		return authorizer.Authorizer{}, err
	}

	valuePtr, gErr := builder.env.GetPointee(returnPtr)
	if gErr != nil {
		slog.Error("authorizerbuilder_buildAuthenticated failed, unable to get return value", slog.Any("err", gErr))
		return authorizer.Authorizer{}, gErr
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

	// authorizerbuilder_addCode returns a 2-slot Result (ptr_or_err, is_err).
	returnPtr, err := builder.env.GetSmallReturnArea()
	if err != nil {
		return err
	}
	defer builder.env.Free(returnPtr, wasm.SmallReturnAreaSize)

	_, err = builder.env.Call(function, returnPtr, builder.ptr, strPtr, uint64(len(code)))
	if err != nil {
		return fmt.Errorf("authorizerbuilder_addCode failed: %w", err)
	}

	// Inspect the Result in the return area to surface parser errors instead of panicking later.
	if _, err := builder.env.GetResult2(returnPtr); err != nil {
		return fmt.Errorf("authorizerbuilder_addCode error: %w", err)
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
