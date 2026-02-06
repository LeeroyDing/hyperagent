package wasm

import (
"context"
"fmt"
"os"

"github.com/tetratelabs/wazero"
"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// WASMManager manages sandboxed WASM modules.
type WASMManager struct {
runtime wazero.Runtime
}

// NewWASMManager creates a new WASMManager.
func NewWASMManager(ctx context.Context) *WASMManager {
r := wazero.NewRuntime(ctx)
wasi_snapshot_preview1.MustInstantiate(ctx, r)
return &WASMManager{runtime: r}
}

// Execute executes a function in a WASM module.
func (m *WASMManager) Execute(ctx context.Context, wasmPath string, functionName string, args ...uint64) ([]uint64, error) {
wasmBytes, err := os.ReadFile(wasmPath)
if err != nil {
return nil, fmt.Errorf("failed to read WASM file: %w", err)
}

mod, err := m.runtime.Instantiate(ctx, wasmBytes)
if err != nil {
return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
}
defer mod.Close(ctx)

f := mod.ExportedFunction(functionName)
if f == nil {
return nil, fmt.Errorf("function %s not found in WASM module", functionName)
}

return f.Call(ctx, args...)
}

// Close closes the WASM runtime.
func (m *WASMManager) Close(ctx context.Context) error {
return m.runtime.Close(ctx)
}
