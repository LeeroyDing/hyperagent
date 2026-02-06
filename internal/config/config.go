package config

import (
"os"
"path/filepath"

"github.com/LeeroyDing/hyperagent/internal/mcp"
"gopkg.in/yaml.v3"
)

type Config struct {
Model            string             `yaml:"model"`
MCPServers       []mcp.ServerConfig `yaml:"mcp_servers"`
InteractiveMode  bool               `yaml:"interactive_mode"`
CommandAllowlist []string           `yaml:"command_allowlist"`
GeminiAPIKey     string             `yaml:"gemini_api_key"`
}

func GetDefaultConfigPath() string {
home, _ := os.UserHomeDir()
return filepath.Join(home, ".hyperagent", "config.yaml")
}

func LoadConfig(path string) (*Config, error) {
if path == "" {
path = GetDefaultConfigPath()
}

f, err := os.Open(path)
if err != nil {
return nil, err
}
defer f.Close()

var cfg Config
decoder := yaml.NewDecoder(f)
err = decoder.Decode(&cfg)
if err != nil {
return nil, err
}

if cfg.Model == "" {
cfg.Model = "gemini-3-flash-preview"
}

return &cfg, nil
}
