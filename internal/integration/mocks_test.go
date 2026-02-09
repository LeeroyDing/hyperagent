package integration

import (
"context"

"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/google/generative-ai-go/genai"
"github.com/philippgille/chromem-go"
)

// MockGeminiClient implements the gemini.GeminiClient interface for testing.
type MockGeminiClient struct {
Responses     []string
ToolCalls     [][]gemini.ToolCall
ResponseIndex int
}

func (m *MockGeminiClient) GenerateContent(ctx context.Context, messages []gemini.Message, tools []*genai.Tool) (string, []gemini.ToolCall, error) {
if m.ResponseIndex >= len(m.Responses) {
return "Mock response", nil, nil
}
response := m.Responses[m.ResponseIndex]
var toolCalls []gemini.ToolCall
if m.ResponseIndex < len(m.ToolCalls) {
toolCalls = m.ToolCalls[m.ResponseIndex]
}
m.ResponseIndex++
return response, toolCalls, nil
}

func (m *MockGeminiClient) SendToolResponse(ctx context.Context, messages []gemini.Message, tools []*genai.Tool, toolResponses []gemini.ToolResponse) (string, []gemini.ToolCall, error) {
if m.ResponseIndex >= len(m.Responses) {
return "Mock tool response", nil, nil
}
response := m.Responses[m.ResponseIndex]
var toolCalls []gemini.ToolCall
if m.ResponseIndex < len(m.ToolCalls) {
toolCalls = m.ToolCalls[m.ResponseIndex]
}
m.ResponseIndex++
return response, toolCalls, nil
}

func (m *MockGeminiClient) EmbedContent(ctx context.Context, text string) ([]float32, error) {
return []float32{0.1, 0.2, 0.3}, nil
}

func (m *MockGeminiClient) Close() error {
return nil
}

type MockEmbedder struct{}

func (m *MockEmbedder) EmbedContent(ctx context.Context, text string) ([]float32, error) {
return []float32{0.1, 0.2, 0.3}, nil
}

type MockExecutor struct {
ExecutedCommands []string
}

func (m *MockExecutor) Execute(sessionID, command string) (string, error) {
m.ExecutedCommands = append(m.ExecutedCommands, command)
return "Mock output for: " + command, nil
}

// MockMemory implements the memory.Memory interface for testing.
type MockMemory struct {
Memorized map[string]string
}

func (m *MockMemory) Memorize(ctx context.Context, id, content string, metadata map[string]string) error {
if m.Memorized == nil {
m.Memorized = make(map[string]string)
}
m.Memorized[id] = content
return nil
}

func (m *MockMemory) Recall(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
return []chromem.Result{}, nil
}

func (m *MockMemory) Forget(ctx context.Context, id string) error {
if m.Memorized != nil {
delete(m.Memorized, id)
}
return nil
}

func (m *MockMemory) Search(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
return m.Recall(ctx, query, limit)
}

func (m *MockMemory) List(ctx context.Context) ([]chromem.Document, error) {
return nil, nil
}

// MockHistory implements the history.History interface for testing.
type MockHistory struct {
Sessions map[string][]history.Message
}

func (h *MockHistory) CreateSession(name string) (string, error) {
id := "mock-session-id"
if h.Sessions == nil {
h.Sessions = make(map[string][]history.Message)
}
h.Sessions[id] = []history.Message{}
return id, nil
}

func (h *MockHistory) AddMessage(sessionID, role, content string) error {
if h.Sessions == nil {
h.Sessions = make(map[string][]history.Message)
}
h.Sessions[sessionID] = append(h.Sessions[sessionID], history.Message{
Role:    role,
Content: content,
})
return nil
}

func (h *MockHistory) LoadHistory(sessionID string) ([]history.Message, error) {
if h.Sessions == nil {
return []history.Message{}, nil
}
return h.Sessions[sessionID], nil
}

func (h *MockHistory) ListSessions() ([]history.Session, error) {
return []history.Session{}, nil
}

func (h *MockHistory) SetSessionName(sessionID, name string) error {
return nil
}

func (h *MockHistory) GetSessionName(sessionID string) string {
return "Mock Session"
}
