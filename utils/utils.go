package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// ResolveAssetFile tries to resolve the .wasm file path dynamically.
//
//	Look for the default relative path from the current working directory and its parents
//	   i.e., walk up until repo root to find target/.../biscuit_wasm_go.wasm
func ResolveAssetFile(path string) (string, error) {

	// Walk up directories to find the default path
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	prev := ""
	dir := cwd
	for dir != prev { // stop when we reach the filesystem root
		candidate := filepath.Join(dir, path)
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, nil
		}
		prev = dir
		dir = filepath.Dir(dir)
	}

	return "", errors.New("wasm file not found; ensure target build exists")
}
