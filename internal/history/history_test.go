package history

import (
"os"
"testing"
)

func TestHistoryManager(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "history_test")
if err != nil {
t.Fatal(err)
}
defer os.RemoveAll(tmpDir)

h, err := NewHistoryManager(tmpDir)
if err != nil {
t.Fatalf("Failed to create HistoryManager: %v", err)
}

sessionID := "test-session"
err = h.AddMessage(sessionID, "user", "hello")
if err != nil {
t.Errorf("AddMessage failed: %v", err)
}

messages, err := h.LoadHistory(sessionID)
if err != nil {
t.Errorf("LoadHistory failed: %v", err)
}

if len(messages) != 1 {
t.Errorf("Expected 1 message, got %d", len(messages))
}

if messages[0].Content != "hello" {
t.Errorf("Expected 'hello', got '%s'", messages[0].Content)
}
}
