package config

import (
"bufio"
"fmt"
"os"
"path/filepath"
"strings"

"gopkg.in/yaml.v3"
)

func RunOOBE() (*Config, error) {
fmt.Println("🚀 Welcome to Hyperagent!")
fmt.Println("It looks like you haven't configured Hyperagent yet.")
fmt.Println("Let's get you set up in a few steps.")
fmt.Println("")

reader := bufio.NewReader(os.Stdin)

// 1. Gemini API Key
fmt.Print("🔑 Enter your Gemini API Key: ")
apiKey, _ := reader.ReadString('\n')
apiKey = strings.TrimSpace(apiKey)

// 2. Model Selection
fmt.Print("🤖 Enter default model [gemini-3-flash-preview]: ")
model, _ := reader.ReadString('\n')
model = strings.TrimSpace(model)
if model == "" {
model = "gemini-3-flash-preview"
}

// 3. Interactive Mode
fmt.Print("🛡️ Enable interactive safety mode by default? (y/n) [y]: ")
interactiveStr, _ := reader.ReadString('\n')
interactiveStr = strings.TrimSpace(strings.ToLower(interactiveStr))
interactive := true
if interactiveStr == "n" {
interactive = false
}

cfg := &Config{
GeminiAPIKey:    apiKey,
Model:           model,
InteractiveMode: interactive,
CommandAllowlist: []string{"ls", "pwd", "cat", "grep", "find"},
}

path := GetDefaultConfigPath()
if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
return nil, fmt.Errorf("failed to create config directory: %w", err)
}

f, err := os.Create(path)
if err != nil {
return nil, fmt.Errorf("failed to create config file: %w", err)
}
defer f.Close()

encoder := yaml.NewEncoder(f)
if err := encoder.Encode(cfg); err != nil {
return nil, fmt.Errorf("failed to save config: %w", err)
}

fmt.Println("")
fmt.Printf("✅ Configuration saved to %s\n", path)
fmt.Println("You're all set! Starting Hyperagent...")

return cfg, nil
}
