package main

import (
"fmt"
"io"
"net/http"
"os"
"os/exec"
"path/filepath"
"testing"
"time"

"github.com/stretchr/testify/assert"
)

func TestStartupIntegration(t *testing.T) {
tmpDir := t.TempDir()
binPath := filepath.Join(tmpDir, "hyperagent")

// Build binary
buildCmd := exec.Command("go", "build", "-o", binPath, "main.go")
if err := buildCmd.Run(); err != nil {
t.Fatalf("failed to build binary: %v", err)
}

// Create a dummy config for testing
testConfig := filepath.Join(tmpDir, "test_config.yaml")
os.WriteFile(testConfig, []byte("gemini_api_key: test-key\nmodel: gemini-3-flash-preview"), 0644)

t.Run("VersionFlag", func(t *testing.T) {
cmd := exec.Command(binPath, "--version")
out, err := cmd.CombinedOutput()
assert.NoError(t, err)
assert.Contains(t, string(out), "v0.0.")
})

t.Run("WebInterfaceE2E", func(t *testing.T) {
port := 3012
cmd := exec.Command(binPath, "--web", "--port", fmt.Sprintf("%d", port), "--config", testConfig)

logPath := filepath.Join(tmpDir, "web_test.log")
logFile, _ := os.Create(logPath)
cmd.Stdout = logFile
cmd.Stderr = logFile

if err := cmd.Start(); err != nil {
t.Fatalf("failed to start web server: %v", err)
}
defer cmd.Process.Kill()

success := false
var lastErr error
for i := 0; i < 20; i++ {
time.Sleep(500 * time.Millisecond)
resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/ui/", port))
if err == nil {
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
if resp.StatusCode == http.StatusOK && assert.Contains(t, string(body), "Hyperagent Web UI") {
success = true
break
}
}
lastErr = err
}

if !success {
logs, _ := os.ReadFile(logPath)
t.Errorf("web server failed to serve UI at port %d: %v\nLogs:\n%s", port, lastErr, string(logs))
}
})
}
