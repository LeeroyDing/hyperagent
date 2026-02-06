package memory

import (
"context"
"fmt"
"runtime"

"github.com/philippgille/chromem-go"
)

// Embedder defines the interface for generating embeddings.
type Embedder interface {
EmbedContent(ctx context.Context, text string) ([]float32, error)
}

type Memory struct {
db         *chromem.DB
collection *chromem.Collection
embedder   Embedder
}

func NewMemory(ctx context.Context, embedder Embedder) (*Memory, error) {
db := chromem.NewDB()
collection, err := db.CreateCollection("agent_memory", nil, nil)
if err != nil {
return nil, fmt.Errorf("failed to create collection: %w", err)
}

return &Memory{
db:         db,
collection: collection,
embedder:   embedder,
}, nil
}

func (m *Memory) Memorize(ctx context.Context, id, content string, metadata map[string]string) error {
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

func (m *Memory) Recall(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
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

func (m *Memory) Forget(ctx context.Context, id string) error {
return m.collection.Delete(ctx, nil, nil, id)
}
