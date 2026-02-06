package agent

import (
"context"
"testing"

"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/google/generative-ai-go/genai"
)

type mockGemini struct{}

func (m *mockGemini) GenerateContent(ctx context.Context, messages []gemini.Message, tools []*genai.Tool) (string, []gemini.ToolCall, error) {
return "This is a distilled summary.", nil, nil
}

func (m *mockGemini) SendToolResponse(ctx context.Context, messages []gemini.Message, tools []*genai.Tool, toolResponses []gemini.ToolResponse) (string, []gemini.ToolCall, error) {
return "", nil, nil
}

type mockEmbedder struct{}

func (m *mockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
return []float32{0.1, 0.2, 0.3}, nil
}

func TestDistill(t *testing.T) {
ctx := context.Background()
tmpDir := t.TempDir()

histMgr, _ := history.NewHistoryManager(tmpDir)
sessionID := "test-session"

// Add 6 messages to trigger distillation (> 5)
for i := 0; i < 6; i++ {
histMgr.AddMessage(sessionID, "user", "hello")
}

mem, _ := memory.NewMemory(ctx, &mockEmbedder{})

a := &Agent{
Gemini:  &mockGemini{},
History: histMgr,
Memory:  mem,
}

err := a.Distill(ctx, sessionID)
if err != nil {
t.Errorf("unexpected error: %v", err)
}

// Verify memory contains the distillation
results, _ := mem.Recall(ctx, "distilled summary", 1)
if len(results) == 0 {
t.Error("expected distillation summary in memory")
}
}
