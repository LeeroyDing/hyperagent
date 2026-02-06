package main

import (
"os/exec"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
)

func TestStartupIntegration(t *testing.T) {
// Build the binary first
tmpDir := t.TempDir()
binPath := filepath.Join(tmpDir, "hyperagent")
buildCmd := exec.Command("go", "build", "-o", binPath, "main.go")
buildCmd.Dir = "."
if err := buildCmd.Run(); err != nil {
t.Fatalf("failed to build binary: %v", err)
}

t.Run("MissingConfigTriggerOOBE", func(t *testing.T) {
// Run with a non-existent config path
// We expect it to try to run OOBE and print the welcome message
cmd := exec.Command(binPath, "--config", filepath.Join(tmpDir, "non-existent.yaml"))
out, _ := cmd.CombinedOutput()

assert.Contains(t, string(out), "Welcome to Hyperagent")
})

t.Run("VersionFlag", func(t *testing.T) {
cmd := exec.Command(binPath, "--version")
out, err := cmd.CombinedOutput()
assert.NoError(t, err)
assert.Contains(t, string(out), "Hyperagent")
assert.Contains(t, string(out), "v0.0.")
})
}
