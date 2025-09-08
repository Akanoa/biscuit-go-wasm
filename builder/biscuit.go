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
	return "biscuit_builder_to_string"
}

// New creates a new BiscuitBuilder using a Wasm environment.
// Returns the initialized BiscuitBuilder or an error in case of failure.
func (builder BiscuitBuilder) New(env wasm.WasmEnv) (BiscuitBuilder, error) {
	returnArea, err := env.GetReturnArea()
	if err != nil {
		return BiscuitBuilder{}, err
	}

	_, err = env.Call("biscuit_builder_new", returnArea)
	if err != nil {
		slog.Error("biscuitbuilder_new failed", slog.Any("err", err))
		return BiscuitBuilder{}, err
	}

	ptr, err := env.ResultPointer(returnArea)
	if err != nil {
		return BiscuitBuilder{}, err
	}

	builder.env = env
	builder.ptr = ptr

	return builder, nil
}

// Build builds the Biscuit using the provided private key.
func (builder *BiscuitBuilder) Build(privateKey keypair.PrivateKey) (token.Biscuit, error) {
	if builder.ptr == 0 {
		return token.Biscuit{}, fmt.Errorf("builder not initialized")
	}

	returnArea, err := builder.env.GetReturnArea()
	if err != nil {
		return token.Biscuit{}, err
	}

	_, err = builder.env.Call("biscuit_builder_build_with_private_key", returnArea, builder.ptr, privateKey.Ptr())
	if err != nil {
		slog.Error("biscuitbuilder_build failed", slog.Any("err", err))
		return token.Biscuit{}, err
	}

	ptr, err := builder.env.ResultPointer(returnArea)
	if err != nil {
	}

	return token.Biscuit{}.New(builder.env, ptr), nil
}

// AddCode adds the provided code to the BiscuitBuilder.
func (builder *BiscuitBuilder) AddCode(code string) error {
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

	_, err = builder.env.Call("biscuit_builder_add_code", returnArea, builder.ptr, strPtr, uint64(len(code)))
	if err != nil {
		return fmt.Errorf("biscuitbuilder_addCode failed: %w", err)
	}

	_, err = builder.env.ResultPointer(returnArea)
	return err
}

func (builder BiscuitBuilder) SetRootKeyId(keyId uint64) error {
	if builder.ptr == 0 {
		return fmt.Errorf("builder builder not initialized")
	}

	returnArea, err := builder.env.GetReturnArea()
	if err != nil {
		return err
	}

	_, err = builder.env.Call("biscuit_builder_set_root_key_id", returnArea, builder.ptr, keyId)
	if err != nil {
		return fmt.Errorf("biscuitbuilder_setRootKeyId failed: %w", err)
	}

	_, err = builder.env.ResultPointer(returnArea)
	return err
}

// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
func (builder BiscuitBuilder) ToString() (string, error) {
	if builder.ptr == 0 {
		return "", fmt.Errorf("token builder not initialized")
	}

	return shared.AsString(builder.env, builder)
}
