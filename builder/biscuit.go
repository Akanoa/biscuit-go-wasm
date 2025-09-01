package builder

import (
	"biscuit-wasm-go/keypair"
	"biscuit-wasm-go/shared"
	"biscuit-wasm-go/token"
	"biscuit-wasm-go/wasm"
	"fmt"
	"log/slog"
)

type BiscuitBuilder struct {
	env wasm.WasmEnv
	ptr uint64
}

func (builder BiscuitBuilder) Ptr() uint64 {
	return builder.ptr
}

func (builder BiscuitBuilder) ToStringWasmFunction() string {
	return "biscuitbuilder_toString"
}

// New creates a new BiscuitBuilder using a Wasm environment.
// Returns the initialized BiscuitBuilder or an error in case of failure.
func (builder BiscuitBuilder) New(env wasm.WasmEnv) (BiscuitBuilder, error) {
	function, err := env.GetFunction("biscuitbuilder_new")
	if err != nil {
		return BiscuitBuilder{}, err
	}

	ret, err := env.Call(function)
	if err != nil {
		slog.Error("biscuitbuilder_new failed", slog.Any("err", err))
		return BiscuitBuilder{}, err
	}

	builder.env = env
	builder.ptr = ret[0]

	return builder, nil
}

// Build builds the Biscuit using the provided private key.
func (builder *BiscuitBuilder) Build(privateKey keypair.PrivateKey) (token.Biscuit, error) {
	if builder.ptr == 0 {
		return token.Biscuit{}, fmt.Errorf("builder not initialized")
	}

	function, err := builder.env.GetFunction("biscuitbuilder_build")
	if err != nil {
		return token.Biscuit{}, err
	}

	returnPtr, err := builder.env.GetReturnArea()
	defer builder.env.Free(returnPtr, wasm.ReturnAreaSize)

	_, err = builder.env.Call(function, returnPtr, builder.ptr, privateKey.Ptr())
	if err != nil {
		slog.Error("biscuitbuilder_build failed", slog.Any("err", err))
		return token.Biscuit{}, err
	}

	valuePtr, err := builder.env.GetPointee(returnPtr)
	if err != nil {
		slog.Error("biscuitbuilder_build failed, unable to get return value", slog.Any("err", err))
		return token.Biscuit{}, err
	}

	return token.Biscuit{}.New(builder.env, valuePtr), nil
}

// AddCode adds the provided code to the BiscuitBuilder.
func (builder *BiscuitBuilder) AddCode(code string) error {
	if builder.ptr == 0 {
		return fmt.Errorf("builder builder not initialized")
	}

	function, err := builder.env.GetFunction("biscuitbuilder_addCode")
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
		return fmt.Errorf("biscuitbuilder_addCode failed: %w", err)
	}
	return nil
}

func (builder BiscuitBuilder) SetRootKeyId(keyId uint64) error {
	if builder.ptr == 0 {
		return fmt.Errorf("builder builder not initialized")
	}

	function, err := builder.env.GetFunction("biscuitbuilder_setRootKeyId")
	if err != nil {
		return err
	}

	_, err = builder.env.Call(function, builder.ptr, keyId)
	if err != nil {
		return fmt.Errorf("biscuitbuilder_setRootKeyId failed: %w", err)
	}
	return nil
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (builder BiscuitBuilder) ToString() (string, error) {
	if builder.ptr == 0 {
		return "", fmt.Errorf("token builder not initialized")
	}

	return shared.AsString(builder.env, builder)
}
