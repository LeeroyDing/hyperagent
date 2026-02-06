package config

import (
"os"
"gopkg.in/yaml.v3"
"github.com/LeeroyDing/hyperagent/internal/mcp"
)

type Config struct {
Model            string             `yaml:"model"` 
MCPServers       []mcp.ServerConfig `yaml:"mcp_servers"` 
InteractiveMode  bool               `yaml:"interactive_mode"` 
CommandAllowlist []string           `yaml:"command_allowlist"` 
}

func LoadConfig(path string) (*Config, error) {
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
