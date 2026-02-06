package config

import (
"os"
"testing"
)

func TestLoadConfig(t *testing.T) {
content := `
interactive_mode: true
command_allowlist:
  - ls
  - echo
mcp_servers:
  - name: test-server
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
`
tmpfile, err := os.CreateTemp("", "config_test.yaml")
if err != nil {
t.Fatal(err)
}
defer os.Remove(tmpfile.Name())

if _, err := tmpfile.Write([]byte(content)); err != nil {
t.Fatal(err)
}
if err := tmpfile.Close(); err != nil {
t.Fatal(err)
}

cfg, err := LoadConfig(tmpfile.Name())
if err != nil {
t.Fatalf("LoadConfig() error = %v", err)
}

if !cfg.InteractiveMode {
t.Errorf("LoadConfig() InteractiveMode = %v, want true", cfg.InteractiveMode)
}

if len(cfg.CommandAllowlist) != 2 || cfg.CommandAllowlist[0] != "ls" {
t.Errorf("LoadConfig() CommandAllowlist = %v, want [ls echo]", cfg.CommandAllowlist)
}
}
