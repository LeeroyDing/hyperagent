package history

import (
"os"
"path/filepath"
"testing"
"github.com/stretchr/testify/assert"
)

func TestFileHistory_Errors(t *testing.T) {
tmpDir, _ := os.MkdirTemp("", "history-err-*")
defer os.RemoveAll(tmpDir)

h, _ := NewHistoryManager(tmpDir)

t.Run("AddMessage_FileError", func(t *testing.T) {
// Create a directory where the file should be to cause an error
sessionID := "test-err"
err := os.MkdirAll(h.GetSessionPath(sessionID), 0755)
assert.NoError(t, err)

err = h.AddMessage(sessionID, "user", "hello")
assert.Error(t, err)
})

t.Run("LoadHistory_CorruptedJSON", func(t *testing.T) {
sessionID := "corrupted"
path := h.GetSessionPath(sessionID)
_ = os.WriteFile(path, []byte("invalid json"), 0644)

msgs, err := h.LoadHistory(sessionID)
assert.Error(t, err)
assert.Nil(t, msgs)
})

t.Run("NewHistoryManager_PermissionError", func(t *testing.T) {
// Use a path that cannot be created
_, err := NewHistoryManager("/proc/invalid/path")
assert.Error(t, err)
})

t.Run("GetSessionName_CorruptedMeta", func(t *testing.T) {
sessionID := "badmeta"
path := h.GetMetadataPath(sessionID)
_ = os.WriteFile(path, []byte("invalid json"), 0644)

name := h.GetSessionName(sessionID)
assert.Equal(t, "New Conversation", name)
})
}
