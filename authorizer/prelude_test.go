package authorizer_test

import (
	"biscuit-wasm-go/wasm"
	"log"
)

var env wasm.WasmEnv

func init() {
	log.Println("setupSuite")
	result, err := wasm.InitWasm()
	if err != nil {
		log.Fatal(err)
		return
	}
	env = result
}
