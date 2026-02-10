package config

import (
"os"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
)

func TestRunOOBE(t *testing.T) {
// Create a temporary home directory
tmpHome, err := os.MkdirTemp("", "hyperagent-home-*")
assert.NoError(t, err)
defer os.RemoveAll(tmpHome)

// Mock HOME environment variable
oldHome := os.Getenv("HOME")
os.Setenv("HOME", tmpHome)
defer os.Setenv("HOME", oldHome)

// Mock Stdin
r, w, err := os.Pipe()
assert.NoError(t, err)
oldStdin := os.Stdin
os.Stdin = r
defer func() {
os.Stdin = oldStdin
}()

// Provide inputs: API Key, Model (default), Interactive (default)
go func() {
defer w.Close()
w.Write([]byte("test-api-key\n"))
w.Write([]byte("\n")) // Default model
w.Write([]byte("\n")) // Default interactive (y)
}()

cfg, err := RunOOBE()
assert.NoError(t, err)
assert.NotNil(t, cfg)

assert.Equal(t, "test-api-key", cfg.GeminiAPIKey)
assert.Equal(t, "gemini-3-flash-preview", cfg.Model)
assert.True(t, cfg.InteractiveMode)

// Verify file was created
configPath := filepath.Join(tmpHome, ".hyperagent", "config.yaml")
_, err = os.Stat(configPath)
assert.NoError(t, err)
}

func TestRunOOBE_CustomInputs(t *testing.T) {
// Create a temporary home directory
tmpHome, err := os.MkdirTemp("", "hyperagent-home-custom-*")
assert.NoError(t, err)
defer os.RemoveAll(tmpHome)

// Mock HOME environment variable
oldHome := os.Getenv("HOME")
os.Setenv("HOME", tmpHome)
defer os.Setenv("HOME", oldHome)

// Mock Stdin
r, w, err := os.Pipe()
assert.NoError(t, err)
oldStdin := os.Stdin
os.Stdin = r
defer func() {
os.Stdin = oldStdin
}()

// Provide inputs: API Key, Custom Model, Interactive (n)
go func() {
defer w.Close()
w.Write([]byte("custom-api-key\n"))
w.Write([]byte("custom-model\n"))
w.Write([]byte("n\n"))
}()

cfg, err := RunOOBE()
assert.NoError(t, err)
assert.NotNil(t, cfg)

assert.Equal(t, "custom-api-key", cfg.GeminiAPIKey)
assert.Equal(t, "custom-model", cfg.Model)
assert.False(t, cfg.InteractiveMode)
}
