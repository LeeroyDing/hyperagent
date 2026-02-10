package config

import (
"os"
"testing"

"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfigPath(t *testing.T) {
path := GetDefaultConfigPath()
assert.NotEmpty(t, path)
assert.Contains(t, path, ".hyperagent")
assert.Contains(t, path, "config.yaml")
}

func TestLoadConfig(t *testing.T) {
t.Run("Success", func(t *testing.T) {
content := `
model: custom-model
interactive_mode: true
command_allowlist:
  - ls
`
tmpfile, err := os.CreateTemp("", "config_success.yaml")
assert.NoError(t, err)
defer os.Remove(tmpfile.Name())
err = os.WriteFile(tmpfile.Name(), []byte(content), 0644)
assert.NoError(t, err)

cfg, err := LoadConfig(tmpfile.Name())
assert.NoError(t, err)
assert.Equal(t, "custom-model", cfg.Model)
assert.True(t, cfg.InteractiveMode)
assert.Equal(t, []string{"ls"}, cfg.CommandAllowlist)
})

t.Run("DefaultModel", func(t *testing.T) {
content := "interactive_mode: false"
tmpfile, err := os.CreateTemp("", "config_default.yaml")
assert.NoError(t, err)
defer os.Remove(tmpfile.Name())
err = os.WriteFile(tmpfile.Name(), []byte(content), 0644)
assert.NoError(t, err)

cfg, err := LoadConfig(tmpfile.Name())
assert.NoError(t, err)
assert.Equal(t, "gemini-3-flash-preview", cfg.Model)
})

t.Run("FileNotFound", func(t *testing.T) {
_, err := LoadConfig("non_existent_file.yaml")
assert.Error(t, err)
})

t.Run("InvalidYAML", func(t *testing.T) {
tmpfile, err := os.CreateTemp("", "config_invalid.yaml")
assert.NoError(t, err)
defer os.Remove(tmpfile.Name())
err = os.WriteFile(tmpfile.Name(), []byte("invalid: yaml: ["), 0644)
assert.NoError(t, err)

_, err = LoadConfig(tmpfile.Name())
assert.Error(t, err)
})

t.Run("DefaultPath", func(t *testing.T) {
// Trigger the branch in LoadConfig(path == "")
_, _ = LoadConfig("")
})
}
