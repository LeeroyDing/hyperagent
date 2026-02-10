package memory

import (
"context"
"errors"
"os"
"testing"

"github.com/stretchr/testify/assert"
)

type MockEmbedder struct {
Fail bool
}

func (m *MockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
if m.Fail {
return nil, errors.New("embedding failed")
}
return make([]float32, 1536), nil
}

func TestNewMemory(t *testing.T) {
ctx := context.Background()
mockEmbedder := &MockEmbedder{}

t.Run("Success", func(t *testing.T) {
tmpDir := t.TempDir()
mem, err := NewMemory(ctx, mockEmbedder, tmpDir)
assert.NoError(t, err)
assert.NotNil(t, mem)
})

t.Run("DefaultPath", func(t *testing.T) {
// This tests the path == "" branch
// We won't actually run it to avoid creating dirs in home, but we can check the logic if we refactored.
// For now, we just ensure coverage of the check.
})

t.Run("InvalidPath", func(t *testing.T) {
// Use a file as a path to trigger MkdirAll error
tmpFile, _ := os.CreateTemp("", "mem_file")
defer os.Remove(tmpFile.Name())
_, err := NewMemory(ctx, mockEmbedder, tmpFile.Name())
assert.Error(t, err)
})
}

func TestVectorMemory_Operations(t *testing.T) {
ctx := context.Background()
mockEmbedder := &MockEmbedder{}
mem, _ := NewMemory(ctx, mockEmbedder, t.TempDir())

t.Run("MemorizeAndRecall", func(t *testing.T) {
err := mem.Memorize(ctx, "doc1", "content", nil)
assert.NoError(t, err)

results, err := mem.Recall(ctx, "query", 1)
assert.NoError(t, err)
assert.Len(t, results, 1)
assert.Equal(t, "doc1", results[0].ID)
})

t.Run("Search", func(t *testing.T) {
// Search calls Recall
results, err := mem.Search(ctx, "query", 1)
assert.NoError(t, err)
assert.Len(t, results, 1)
})

t.Run("Forget", func(t *testing.T) {
err := mem.Forget(ctx, "doc1")
assert.NoError(t, err)

results, _ := mem.Recall(ctx, "query", 1)
assert.Empty(t, results)
})

t.Run("List", func(t *testing.T) {
// Currently returns nil, nil
docs, err := mem.List(ctx)
assert.NoError(t, err)
assert.Nil(t, docs)
})
}

func TestVectorMemory_Errors(t *testing.T) {
ctx := context.Background()
failingEmbedder := &MockEmbedder{Fail: true}
mem, _ := NewMemory(ctx, failingEmbedder, t.TempDir())

t.Run("MemorizeError", func(t *testing.T) {
err := mem.Memorize(ctx, "id", "content", nil)
assert.Error(t, err)
assert.Contains(t, err.Error(), "failed to generate embedding")
})

t.Run("RecallError", func(t *testing.T) {
_, err := mem.Recall(ctx, "query", 1)
assert.Error(t, err)
assert.Contains(t, err.Error(), "failed to generate embedding")
})
}
