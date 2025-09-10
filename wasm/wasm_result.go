// Provides facilities for reading Wasm results.

package wasm

import (
	error2 "biscuit-wasm-go/error"
	"encoding/binary"
	"fmt"
)

// WasmResult is the result of a Wasm function call.
// It is a 12-byte struct with the following layout:
// Layout:
// 0: 4 bytes: data pointer
// 4: 4 length: data length
// 8: 4 bytes: is_ok (1 = Ok, non-zero = Err)
const wasmResultSize = 12

type Result struct {
	data uint32
	ln   uint32
	kind ResultKind
	env  WasmEnv
}

type ResultKind int

const (
	Ok            ResultKind = 0
	Biscuit                  = 1
	Serialization            = 2
)

func toBisuitError(data []byte, kind ResultKind) error2.BiscuitError {

	if kind == Serialization {
		return error2.BiscuitError{
			Raw: string(data),
		}
	}

	biscuitError, err := error2.BiscuitError{}.FromString(string(data))
	if err != nil {
		return error2.BiscuitError{
			Raw: string(data),
		}
	}
	return biscuitError

}

func (kind ResultKind) isError() bool {
	if kind != Ok {
		return true
	}
	return false
}

// readWasmResult reads the result from the memory.
func (env WasmEnv) readWasmResult(ptr uint64) (Result, error) {
	mem := env.Module.Memory()
	buf, ok := mem.Read(uint32(ptr), wasmResultSize)
	if !ok {
		return Result{}, fmt.Errorf("failed to read result at addr %x", ptr)
	}

	dataPtr := binary.LittleEndian.Uint32(buf[0:4])
	dataLn := binary.LittleEndian.Uint32(buf[4:8])
	kind := binary.LittleEndian.Uint32(buf[8:12])

	return Result{
		data: dataPtr,
		ln:   dataLn,
		kind: ResultKind(kind),
		env:  env,
	}, nil
}

// readDataBytes reads the data bytes from the memory.
func (result Result) readDataBytes() ([]byte, error) {
	mem := result.env.Module.Memory()
	buf, ok := mem.Read(result.data, result.ln)
	if !ok {
		return nil, fmt.Errorf("failed to read data result at addr %x", result.data)
	}
	return buf, nil
}

// asPtr returns the pointer to the result or the string error message if the result is an error.
// Return the pointer to the result if the result is Ok.
// Otherwise, read the memory to get the error message string.
func (result Result) asPtr() (uint64, error) {
	if result.kind == Ok {
		return uint64(result.data), nil
	}

	errString, err := result.readDataBytes()
	if err != nil {
		return 0, err
	}

	return 0, toBisuitError(errString, result.kind)

}

// asNumber returns the number result or the string error message if the result is an error.
// Return the number result if the result is Ok.
// Otherwise, read the memory to get the error message string.
func (result Result) asNumber() (uint32, error) {
	if result.kind == Ok {
		return result.data, nil
	}

	errString, err := result.readDataBytes()
	if err != nil {
		return 0, err
	}

	return 0, toBisuitError(errString, result.kind)
}

// asString returns the string result or the string error message if the result is an error.
// Read the memory to get the string data.
// Return the string result if the result is Ok.
// Otherwise, the string is the error message.
func (result Result) asString() (string, error) {
	data, err := result.readDataBytes()

	if err != nil {
		return "", err
	}

	if result.kind == Ok {
		return string(data), nil
	}

	return "", toBisuitError(data, result.kind)
}

// asBytes returns bytes result or the string error message if the result is an error.
// Read the memory to get the bytes data.
// Return the bytes result if the result is Ok.
// Otherwise, encode data as UTF-8 and return the string.
//
// Memory Layout Diagram:
// +----------------+     +-------------------+
// | Return Area    |     | Bytes Data        |
// | (8 bytes)      |     | (variable length) |
// +----------------+     +-------------------+
// | Bytes Ptr    --|---->| Actual bytes    |
// | Bytes Length   |     | content...       |
// +----------------+     +-------------------+
//
//	^
//	|
//
// ptr (input parameter)
func (result Result) asBytes() ([]byte, error) {
	data, err := result.readDataBytes()

	if err != nil {
		return nil, err
	}

	if result.kind == Ok {
		return data, nil
	}

	return nil, toBisuitError(data, result.kind)
}

//---------------- Public functions  ------------------

// ResultPointer returns the pointer to the result or the string error message if the result is an error.
// Result<Box<T>, String>
// Pointer is at user charge to be freed.
func (env WasmEnv) ResultPointer(ptr uint64) (uint64, error) {
	result, err := env.readWasmResult(ptr)
	if err != nil {
		return 0, err
	}
	return result.asPtr()
}

// ResultString returns the string result or the string error message if the result is an error.
// Result<String, String>
// Pointer is at user charge to be freed.
func (env WasmEnv) ResultString(ptr uint64) (string, error) {
	result, err := env.readWasmResult(ptr)
	if err != nil {
		return "", err
	}
	return result.asString()
}

// ResultBytes returns the bytes result or the string error message if the result is an error.
// Result<Vec<u8>, String>
// Pointer is at user charge to be freed.
func (env WasmEnv) ResultBytes(ptr uint64) ([]byte, error) {
	result, err := env.readWasmResult(ptr)
	if err != nil {
		return nil, err
	}
	return result.asBytes()
}

// ResultNumber returns the number result or the string error message if the result is an error.
// Result<u32, String>
// Pointer is at user charge to be freed.
func (env WasmEnv) ResultNumber(ptr uint64) (uint32, error) {
	result, err := env.readWasmResult(ptr)
	if err != nil {
		return 0, err
	}
	return result.asNumber()
}
