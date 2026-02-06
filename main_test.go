package main

import (
"testing"
"github.com/LeeroyDing/hyperagent/internal/executor"
)

func TestNewShellExecutor(t *testing.T) {
// NewShellExecutor now takes a slice of strings for the allowlist
exec := executor.NewShellExecutor([]string{"ls", "echo"})
if exec == nil {
t.Fatal("Expected non-nil executor")
}
}
