package web

import (
"bytes"
"net/http"
"net/http/httptest"
"os"
"path/filepath"
"testing"

"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/gin-gonic/gin"
)

func TestConfigEndpoints(t *testing.T) {
gin.SetMode(gin.TestMode)
tmpDir := t.TempDir()
configPath := filepath.Join(tmpDir, "config.yaml")
initialConfig := "gemini_api_key: test-key"
os.WriteFile(configPath, []byte(initialConfig), 0644)

histMgr, _ := history.NewHistoryManager(filepath.Join(tmpDir, "history"))
s := NewServer(nil, histMgr)
s.ConfigPath = configPath

// Test GET /api/config
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/config", nil)
s.Router.ServeHTTP(w, req)

if w.Code != http.StatusOK {
t.Errorf("expected 200, got %d", w.Code)
}
if w.Body.String() != initialConfig {
t.Errorf("expected %s, got %s", initialConfig, w.Body.String())
}

// Test POST /api/config
newConfig := "gemini_api_key: updated-key"
w = httptest.NewRecorder()
req, _ = http.NewRequest("POST", "/api/config", bytes.NewBufferString(newConfig))
s.Router.ServeHTTP(w, req)

if w.Code != http.StatusOK {
t.Errorf("expected 200, got %d", w.Code)
}

// Verify file was updated
data, _ := os.ReadFile(configPath)
if string(data) != newConfig {
t.Errorf("expected %s, got %s", newConfig, string(data))
}
}
