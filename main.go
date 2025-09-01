package main

import (
	"biscuit-wasm-go/wasm"
	"log/slog"
	"os"
)

func main() {

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	_, err := wasm.InitWasm()
	if err != nil {
		panic(err)
	}
}
