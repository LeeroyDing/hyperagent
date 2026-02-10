package memory

import (
"context"
"fmt"
"os"
"path/filepath"
"runtime"

"github.com/philippgille/chromem-go"
)

// Embedder defines the interface for generating embeddings.
type Embedder interface {
EmbedContent(ctx context.Context, text string) ([]float32, error)
}

// Memory defines the interface for long-term memory operations.
type Memory interface {
Memorize(ctx context.Context, id, content string, metadata map[string]string) error
Recall(ctx context.Context, query string, limit int) ([]chromem.Result, error)
Forget(ctx context.Context, id string) error
Search(ctx context.Context, query string, limit int) ([]chromem.Result, error)
List(ctx context.Context) ([]chromem.Document, error)
}

// VectorMemory implements the Memory interface using a vector database.
type VectorMemory struct {
db         *chromem.DB
collection *chromem.Collection
embedder   Embedder
}

// NewMemory creates a new VectorMemory instance.
func NewMemory(ctx context.Context, embedder Embedder, path string) (*VectorMemory, error) {
if path == "" {
home, _ := os.UserHomeDir()
path = filepath.Join(home, ".hyperagent", "memory")
}
if err := os.MkdirAll(path, 0755); err != nil {
return nil, fmt.Errorf("failed to create memory directory: %w", err)
}

// Use persistent storage
db, err := chromem.NewPersistentDB(path, false)
if err != nil {
return nil, fmt.Errorf("failed to create persistent db: %w", err)
}
collection, err := db.GetOrCreateCollection("agent_memory", nil, nil)
if err != nil {
return nil, fmt.Errorf("failed to get or create collection: %w", err)
}

return &VectorMemory{
db:         db,
collection: collection,
embedder:   embedder,
}, nil
}

func (m *VectorMemory) Memorize(ctx context.Context, id, content string, metadata map[string]string) error {
embedding, err := m.embedder.EmbedContent(ctx, content)
if err != nil {
return fmt.Errorf("failed to generate embedding: %w", err)
}

doc := chromem.Document{
ID:        id,
Content:   content,
Metadata:  metadata,
Embedding: embedding,
}

err = m.collection.AddDocuments(ctx, []chromem.Document{doc}, runtime.NumCPU())
if err != nil {
return fmt.Errorf("failed to add document: %w", err)
}

return nil
}

func (m *VectorMemory) Recall(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
embedding, err := m.embedder.EmbedContent(ctx, query)
if err != nil {
return nil, fmt.Errorf("failed to generate embedding for query: %w", err)
}

results, err := m.collection.QueryEmbedding(ctx, embedding, limit, nil, nil)
if err != nil {
return nil, fmt.Errorf("failed to query collection: %w", err)
}

return results, nil
}

func (m *VectorMemory) Forget(ctx context.Context, id string) error {
return m.collection.Delete(ctx, nil, nil, id)
}

func (m *VectorMemory) Search(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
return m.Recall(ctx, query, limit)
}

func (m *VectorMemory) List(ctx context.Context) ([]chromem.Document, error) {
return nil, nil // Not implemented yet
}
