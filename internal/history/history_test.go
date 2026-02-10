package history

import (
"os"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
)

func TestNewHistoryManager(t *testing.T) {
t.Run("CustomDir", func(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "history_test_custom")
assert.NoError(t, err)
defer os.RemoveAll(tmpDir)

h, err := NewHistoryManager(tmpDir)
assert.NoError(t, err)
assert.Equal(t, tmpDir, h.StorageDir)
})

t.Run("DefaultDir", func(t *testing.T) {
path := GetDefaultHistoryDir()
assert.NotEmpty(t, path)
})

t.Run("CreationError", func(t *testing.T) {
tmpFile, _ := os.CreateTemp("", "not_a_dir")
defer os.Remove(tmpFile.Name())
_, err := NewHistoryManager(tmpFile.Name())
assert.Error(t, err)
})
}

func TestFileHistory_SessionOperations(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "history_test_ops")
assert.NoError(t, err)
defer os.RemoveAll(tmpDir)

h, _ := NewHistoryManager(tmpDir)

t.Run("CreateAndGetName", func(t *testing.T) {
id, err := h.CreateSession("Test Session")
assert.NoError(t, err)
assert.NotEmpty(t, id)
assert.Equal(t, "Test Session", h.GetSessionName(id))
})

t.Run("CreateDefaultName", func(t *testing.T) {
id, err := h.CreateSession("")
assert.NoError(t, err)
assert.Equal(t, "New Conversation", h.GetSessionName(id))
})

t.Run("SetSessionName", func(t *testing.T) {
id, _ := h.CreateSession("Old")
err := h.SetSessionName(id, "New")
assert.NoError(t, err)
assert.Equal(t, "New", h.GetSessionName(id))
})

t.Run("GetNonExistentName", func(t *testing.T) {
assert.Equal(t, "New Conversation", h.GetSessionName("none"))
})

t.Run("GetCorruptedName", func(t *testing.T) {
id := "corrupt-meta"
path := h.GetMetadataPath(id)
os.WriteFile(path, []byte("{invalid json"), 0644)
assert.Equal(t, "New Conversation", h.GetSessionName(id))
})

t.Run("CreateSessionError", func(t *testing.T) {
badDir := filepath.Join(tmpDir, "file_not_dir")
os.WriteFile(badDir, []byte("I am a file"), 0644)
h2 := &FileHistory{StorageDir: badDir}
_, err := h2.CreateSession("fail")
assert.Error(t, err)
})
}

func TestFileHistory_Messages(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "history_test_msgs")
assert.NoError(t, err)
defer os.RemoveAll(tmpDir)

h, _ := NewHistoryManager(tmpDir)
sessionID := "msg-test"

t.Run("AddAndLoad", func(t *testing.T) {
err := h.AddMessage(sessionID, "user", "hello")
assert.NoError(t, err)
err = h.AddMessage(sessionID, "assistant", "hi there")
assert.NoError(t, err)

msgs, err := h.LoadHistory(sessionID)
assert.NoError(t, err)
assert.Len(t, msgs, 2)
assert.Equal(t, "user", msgs[0].Role)
assert.Equal(t, "hello", msgs[0].Content)
})

t.Run("LoadNonExistent", func(t *testing.T) {
msgs, err := h.LoadHistory("ghost")
assert.NoError(t, err)
assert.Empty(t, msgs)
})

t.Run("AddMessageError", func(t *testing.T) {
badSessionID := "bad-session"
badPath := h.GetSessionPath(badSessionID)
os.Mkdir(badPath, 0755)
defer os.RemoveAll(badPath)

err := h.AddMessage(badSessionID, "user", "fail")
assert.Error(t, err)
})

t.Run("LoadHistoryCorrupt", func(t *testing.T) {
id := "corrupt-history"
path := h.GetSessionPath(id)
os.WriteFile(path, []byte("{invalid json\n"), 0644)
_, err := h.LoadHistory(id)
assert.Error(t, err)
})
}

func TestFileHistory_ListSessions(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "history_test_list")
assert.NoError(t, err)
defer os.RemoveAll(tmpDir)

h, _ := NewHistoryManager(tmpDir)

h.CreateSession("Session 1")
h.CreateSession("Session 2")
os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("test"), 0644)

sessions, err := h.ListSessions()
assert.NoError(t, err)
assert.Len(t, sessions, 2)

t.Run("ListSessionsError", func(t *testing.T) {
h2 := &FileHistory{StorageDir: "/non/existent/path"}
_, err := h2.ListSessions()
assert.Error(t, err)
})
}
