package memory

import (
"context"
"os"
"path/filepath"
"testing"

"github.com/stretchr/testify/assert"
)

func TestNewMemory_DefaultPath(t *testing.T) {
ctx := context.Background()
mockEmbedder := &MockEmbedder{}

// Save old HOME and restore after test
oldHome := os.Getenv("HOME")
tmpHome := t.TempDir()
os.Setenv("HOME", tmpHome)
defer os.Setenv("HOME", oldHome)

mem, err := NewMemory(ctx, mockEmbedder, "")
assert.NoError(t, err)
assert.NotNil(t, mem)

expectedPath := filepath.Join(tmpHome, ".hyperagent", "memory")
assert.DirExists(t, expectedPath)
}

func TestNewMemory_MkdirError(t *testing.T) {
ctx := context.Background()
mockEmbedder := &MockEmbedder{}

// Create a file where a directory should be
tmpDir := t.TempDir()
filePath := filepath.Join(tmpDir, "blocked")
os.WriteFile(filePath, []byte("data"), 0644)

// Try to create memory inside that file path
_, err := NewMemory(ctx, mockEmbedder, filepath.Join(filePath, "subdir"))
assert.Error(t, err)
assert.Contains(t, err.Error(), "failed to create memory directory")
}

func TestVectorMemory_List_Placeholder(t *testing.T) {
ctx := context.Background()
mem, _ := NewMemory(ctx, &MockEmbedder{}, t.TempDir())
docs, err := mem.List(ctx)
assert.NoError(t, err)
assert.Nil(t, docs)
}
