package memory

import (
"context"
"testing"
)

type MockEmbedder struct{}

func (m *MockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
return []float32{0.1, 0.2, 0.3}, nil
}

func TestMemory(t *testing.T) {
ctx := context.Background()
mockEmbedder := &MockEmbedder{}

mem, err := NewMemory(ctx, mockEmbedder)
if err != nil {
t.Fatalf("NewMemory() error = %v", err)
}

// Test Memorize
err = mem.Memorize(ctx, "doc1", "hello world", map[string]string{"key": "value"})
if err != nil {
t.Fatalf("Memorize() error = %v", err)
}

// Test Recall
results, err := mem.Recall(ctx, "hello", 1)
if err != nil {
t.Fatalf("Recall() error = %v", err)
}
if len(results) != 1 || results[0].ID != "doc1" {
t.Errorf("Recall() results = %v, want 1 result with ID doc1", results)
}

// Test Forget
err = mem.Forget(ctx, "doc1")
if err != nil {
t.Fatalf("Forget() error = %v", err)
}

// Verify forgotten
results, _ = mem.Recall(ctx, "hello", 1)
if len(results) != 0 {
t.Errorf("Recall() after Forget results = %v, want 0", results)
}
}
