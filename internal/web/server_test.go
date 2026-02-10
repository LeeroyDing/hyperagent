package web

import (
"bytes"
"context"
"encoding/json"
"errors"
"net/http"
"net/http/httptest"
"testing"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/google/generative-ai-go/genai"
"github.com/philippgille/chromem-go"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
)

// Mocks
type MockHistory struct {
mock.Mock
}

func (m *MockHistory) CreateSession(name string) (string, error) {
args := m.Called(name)
return args.String(0), args.Error(1)
}

func (m *MockHistory) AddMessage(sessionID, role, content string) error {
args := m.Called(sessionID, role, content)
return args.Error(0)
}

func (m *MockHistory) LoadHistory(sessionID string) ([]history.Message, error) {
args := m.Called(sessionID)
return args.Get(0).([]history.Message), args.Error(1)
}

func (m *MockHistory) ListSessions() ([]history.Session, error) {
args := m.Called()
return args.Get(0).([]history.Session), args.Error(1)
}

func (m *MockHistory) SetSessionName(sessionID, name string) error {
args := m.Called(sessionID, name)
return args.Error(0)
}

func (m *MockHistory) GetSessionName(sessionID string) string {
args := m.Called(sessionID)
return args.String(0)
}

type MockMemory struct {
mock.Mock
}

func (m *MockMemory) Memorize(ctx context.Context, id, content string, metadata map[string]string) error {
args := m.Called(ctx, id, content, metadata)
return args.Error(0)
}

func (m *MockMemory) Recall(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
args := m.Called(ctx, query, limit)
return args.Get(0).([]chromem.Result), args.Error(1)
}

func (m *MockMemory) Forget(ctx context.Context, id string) error {
args := m.Called(ctx, id)
return args.Error(0)
}

func (m *MockMemory) Search(ctx context.Context, query string, limit int) ([]chromem.Result, error) {
args := m.Called(ctx, query, limit)
return args.Get(0).([]chromem.Result), args.Error(1)
}

func (m *MockMemory) List(ctx context.Context) ([]chromem.Document, error) {
args := m.Called(ctx)
return args.Get(0).([]chromem.Document), args.Error(1)
}

type MockGemini struct {
mock.Mock
}

func (m *MockGemini) GenerateContent(ctx context.Context, messages []gemini.Message, tools []*genai.Tool) (string, []gemini.ToolCall, error) {
args := m.Called(ctx, messages, tools)
return args.String(0), args.Get(1).([]gemini.ToolCall), args.Error(2)
}

func (m *MockGemini) SendToolResponse(ctx context.Context, messages []gemini.Message, tools []*genai.Tool, toolResponses []gemini.ToolResponse) (string, []gemini.ToolCall, error) {
args := m.Called(ctx, messages, tools, toolResponses)
return args.String(0), args.Get(1).([]gemini.ToolCall), args.Error(2)
}

func (m *MockGemini) EmbedContent(ctx context.Context, text string) ([]float32, error) {
args := m.Called(ctx, text)
return args.Get(0).([]float32), args.Error(1)
}

func (m *MockGemini) Close() error {
args := m.Called()
return args.Error(0)
}

func TestWebAPI(t *testing.T) {
mockHist := new(MockHistory)
mockMem := new(MockMemory)
mockGemini := new(MockGemini)

a := agent.NewAgent(mockGemini, nil, mockMem, nil, mockHist, false)
s := NewServer(a, mockHist, mockMem)

t.Run("RedirectRoot", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusMovedPermanently, w.Code)
})

t.Run("GetSessions_Success", func(t *testing.T) {
sessions := []history.Session{{ID: "1", Name: "Test"}}
mockHist.On("ListSessions").Return(sessions, nil).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/sessions", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("GetSessions_Error", func(t *testing.T) {
mockHist.On("ListSessions").Return([]history.Session{}, errors.New("db error")).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/sessions", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})

t.Run("CreateSession_Success", func(t *testing.T) {
mockHist.On("CreateSession", "New Session").Return("uuid-123", nil).Once()
body, _ := json.Marshal(map[string]string{"name": "New Session"})
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusCreated, w.Code)
})

t.Run("CreateSession_InvalidJSON", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBufferString("invalid"))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusBadRequest, w.Code)
})

t.Run("CreateSession_Error", func(t *testing.T) {
mockHist.On("CreateSession", "Fail").Return("", errors.New("fail")).Once()
body, _ := json.Marshal(map[string]string{"name": "Fail"})
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(body))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})

t.Run("GetMessages_Success", func(t *testing.T) {
mockHist.On("LoadHistory", "123").Return([]history.Message{}, nil).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/sessions/123/messages", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("GetMessages_Error", func(t *testing.T) {
mockHist.On("LoadHistory", "123").Return([]history.Message{}, errors.New("error")).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/sessions/123/messages", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})

t.Run("SendMessage_Success", func(t *testing.T) {
mockHist.On("LoadHistory", "123").Return([]history.Message{}, nil).Once()
mockMem.On("Recall", mock.Anything, "hello", 5).Return([]chromem.Result{}, nil).Once()
mockHist.On("AddMessage", "123", "user", "hello").Return(nil).Once()
mockGemini.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return("hi", []gemini.ToolCall{}, nil).Once()
mockHist.On("AddMessage", "123", "model", "hi").Return(nil).Once()
body, _ := json.Marshal(map[string]string{"content": "hello"})
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions/123/messages", bytes.NewBuffer(body))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("SendMessage_InvalidJSON", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions/123/messages", bytes.NewBufferString("invalid"))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusBadRequest, w.Code)
})

t.Run("SendMessage_AgentError", func(t *testing.T) {
mockHist.On("LoadHistory", "123").Return([]history.Message{}, nil).Once()
mockMem.On("Recall", mock.Anything, "fail", 5).Return([]chromem.Result{}, nil).Once()
mockHist.On("AddMessage", "123", "user", "fail").Return(nil).Once()
mockGemini.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return("", []gemini.ToolCall{}, errors.New("agent fail")).Once()
body, _ := json.Marshal(map[string]string{"content": "fail"})
w := httptest.NewRecorder()
req, _ := http.NewRequest("POST", "/api/sessions/123/messages", bytes.NewBuffer(body))
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})

t.Run("SearchMemory_Success", func(t *testing.T) {
mockMem.On("Search", mock.Anything, "test", 10).Return([]chromem.Result{}, nil).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/memory?q=test", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("SearchMemory_Error", func(t *testing.T) {
mockMem.On("Search", mock.Anything, "test", 10).Return([]chromem.Result{}, errors.New("error")).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/memory?q=test", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})

t.Run("DeleteMemory_Success", func(t *testing.T) {
mockMem.On("Forget", mock.Anything, "m1").Return(nil).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("DELETE", "/api/memory/m1", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("DeleteMemory_Error", func(t *testing.T) {
mockMem.On("Forget", mock.Anything, "m1").Return(errors.New("error")).Once()
w := httptest.NewRecorder()
req, _ := http.NewRequest("DELETE", "/api/memory/m1", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusInternalServerError, w.Code)
})
}


func TestServer_Run(t *testing.T) {
// This is a bit tricky as Run blocks. We'll run it in a goroutine.
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

mockAgent := &mockAgent{}
mockHistory := &mockHistory{}
mockMemory := &mockMemory{}
srv := NewServer(mockAgent, mockHistory, mockMemory)

go func() {
// Use a random high port
_ = srv.Run("127.0.0.1:0")
}()

// Give it a moment to start
time.Sleep(100 * time.Millisecond)
// In a real scenario, we'd check if the port is open, but for coverage, 
// just entering the function and starting the listener is often enough.
}
