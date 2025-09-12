package builder

//
//import (
//	"biscuit-wasm-go/shared"
//	"biscuit-wasm-go/wasm"
//	"fmt"
//	"log/slog"
//)
//
//type BlockBuilder struct {
//	env wasm.WasmEnv
//	ptr uint64
//}
//
//func (builder BlockBuilder) Ptr() uint64 {
//	return builder.ptr
//}
//
//func (builder BlockBuilder) ToStringWasmFunction() string {
//	return "blockbuilder_toString"
//}
//
//// New creates a new BlockBuilder using a Wasm environment.
//// Returns the initialized BlockBuilder or an error in case of failure.
//func (builder BlockBuilder) New(env wasm.WasmEnv) (BlockBuilder, error) {
//	function, err := env.GetFunction("blockbuilder_new")
//	if err != nil {
//		return BlockBuilder{}, err
//	}
//
//	ret, err := env.Call(function)
//	if err != nil {
//		slog.Error("blockbuilder_new failed", slog.Any("err", err))
//		return BlockBuilder{}, err
//	}
//
//	builder.env = env
//	builder.ptr = ret[0]
//
//	return builder, nil
//}
//
//// AddCode adds the provided code to the block.
//func (block *BlockBuilder) AddCode(code string) error {
//	if block.ptr == 0 {
//		return fmt.Errorf("block builder not initialized")
//	}
//
//	function, err := block.env.GetFunction("blockbuilder_addCode")
//	if err != nil {
//		return err
//	}
//
//	strPtr, err := block.env.WriteString(code)
//	if err != nil {
//		return fmt.Errorf("cannot write string to wasm memory: %w", err)
//	}
//	defer block.env.Free(strPtr, uint64(len(code)))
//
//	returnPtr, err := block.env.GetReturnArea()
//	if err != nil {
//		return err
//	}
//
//	_, err = block.env.Call(function, returnPtr, block.ptr, strPtr, uint64(len(code)))
//	if err != nil {
//		return fmt.Errorf("blockbuilder_addCode failed: %w", err)
//	}
//	return nil
//}
//
//// ToString converts the PublicKey to its string representation using the linked Wasm environment. Returns the string or an error.
//func (builder BlockBuilder) ToString() (string, error) {
//	if builder.ptr == 0 {
//		return "", fmt.Errorf("blockbuilder not initialized")
//	}
//
//	return shared.AsString(builder.env, builder)
//}
