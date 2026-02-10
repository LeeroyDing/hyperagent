package agent

import (
"context"
"errors"
"fmt"
"os"
"testing"

"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/philippgille/chromem-go"
"github.com/stretchr/testify/assert"
)

func TestNewAgent(t *testing.T) {
g := &MockGeminiClient{}
e := &MockExecutor{}
m := &MockMemory{}
h := &MockHistory{}
a := NewAgent(g, e, m, nil, h, false)

assert.NotNil(t, a)
assert.Equal(t, g, a.Gemini)
assert.Equal(t, e, a.Executor)
assert.Equal(t, m, a.Memory)
assert.Equal(t, h, a.History)
assert.NotNil(t, a.Editor)
assert.NotNil(t, a.Orchestrator)
}

func TestAgent_Run(t *testing.T) {
ctx := context.Background()

t.Run("Success with RAG and Tool Call", func(t *testing.T) {
g := &MockGeminiClient{
Responses: []string{"", "Final answer"},
ToolCalls: [][]gemini.ToolCall{
{{Name: "execute_command", Arguments: map[string]interface{}{"command": "ls"}}},
nil,
},
}
e := &MockExecutor{}
m := &MockMemory{
RecallResults: []chromem.Result{{Content: "memory context"}},
}
h := &MockHistory{}
a := NewAgent(g, e, m, nil, h, false)

resp, err := a.Run(ctx, "s1", "hello")
assert.NoError(t, err)
assert.Equal(t, "Final answer", resp)
assert.Contains(t, e.ExecutedCommands, "ls")
})

t.Run("History Load Error", func(t *testing.T) {
h := &MockHistory{LoadError: errors.New("history error")}
a := NewAgent(nil, nil, &MockMemory{}, nil, h, false)
_, err := a.Run(ctx, "s1", "hello")
assert.Error(t, err)
assert.Contains(t, err.Error(), "history error")
})

t.Run("Gemini Generate Error", func(t *testing.T) {
g := &MockGeminiClient{GenerateError: errors.New("gemini error")}
h := &MockHistory{}
a := NewAgent(g, nil, &MockMemory{}, nil, h, false)
_, err := a.Run(ctx, "s1", "hello")
assert.Error(t, err)
assert.Contains(t, err.Error(), "gemini error")
})

t.Run("Gemini Tool Response Error", func(t *testing.T) {
g := &MockGeminiClient{
Responses: []string{""},
ToolCalls: [][]gemini.ToolCall{
{{Name: "execute_command", Arguments: map[string]interface{}{"command": "ls"}}},
},
SendToolResponseError: errors.New("tool response error"),
}
h := &MockHistory{}
a := NewAgent(g, &MockExecutor{}, &MockMemory{}, nil, h, false)
_, err := a.Run(ctx, "s1", "hello")
assert.Error(t, err)
assert.Contains(t, err.Error(), "tool response error")
})
}

func TestAgent_HandleToolCall(t *testing.T) {
ctx := context.Background()

t.Run("execute_command cancelled", func(t *testing.T) {
// Simulate 'n' input
oldStdin := os.Stdin
r, w, _ := os.Pipe()
os.Stdin = r
fmt.Fprintln(w, "n")
w.Close()
defer func() { os.Stdin = oldStdin }()

a := &Agent{InteractiveMode: true}
tc := gemini.ToolCall{Name: "execute_command", Arguments: map[string]interface{}{"command": "ls"}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "Action cancelled by user", resp)
})

t.Run("read_file success with end", func(t *testing.T) {
a := NewAgent(nil, nil, nil, nil, nil, false)
f, _ := os.CreateTemp("", "testfile")
defer os.Remove(f.Name())
f.WriteString("line1\nline2\nline3")
f.Close()

tc := gemini.ToolCall{Name: "read_file", Arguments: map[string]interface{}{"path": f.Name(), "start": 1.0, "end": 2.0}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "line1\nline2", resp)
})

t.Run("read_file error", func(t *testing.T) {
a := NewAgent(nil, nil, nil, nil, nil, false)
tc := gemini.ToolCall{Name: "read_file", Arguments: map[string]interface{}{"path": "nonexistent", "start": 1.0, "end": 2.0}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.Error(t, err)
assert.Empty(t, resp)
})

t.Run("replace_text success", func(t *testing.T) {
a := NewAgent(nil, nil, nil, nil, nil, false)
f, _ := os.CreateTemp("", "testfile")
defer os.Remove(f.Name())
f.WriteString("old")
f.Close()

tc := gemini.ToolCall{Name: "replace_text", Arguments: map[string]interface{}{"path": f.Name(), "old_text": "old", "new_text": "new"}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "Text replaced successfully", resp)
})

t.Run("replace_text cancelled", func(t *testing.T) {
// Simulate 'n' input
oldStdin := os.Stdin
r, w, _ := os.Pipe()
os.Stdin = r
fmt.Fprintln(w, "n")
w.Close()
defer func() { os.Stdin = oldStdin }()

a := &Agent{InteractiveMode: true}
tc := gemini.ToolCall{Name: "replace_text", Arguments: map[string]interface{}{"path": "p", "old_text": "o", "new_text": "n"}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "Action cancelled by user", resp)
})

t.Run("memory_save success", func(t *testing.T) {
m := &MockMemory{}
a := NewAgent(nil, nil, m, nil, nil, false)
tc := gemini.ToolCall{Name: "memory_save", Arguments: map[string]interface{}{"id": "id", "content": "c"}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "Information memorized", resp)
})

t.Run("memory_load success with limit", func(t *testing.T) {
m := &MockMemory{RecallResults: []chromem.Result{{ID: "id", Content: "c"}}}
a := NewAgent(nil, nil, m, nil, nil, false)
tc := gemini.ToolCall{Name: "memory_load", Arguments: map[string]interface{}{"query": "q", "limit": 10.0}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Contains(t, resp, "c")
})

t.Run("memory_forget success", func(t *testing.T) {
m := &MockMemory{}
a := NewAgent(nil, nil, m, nil, nil, false)
tc := gemini.ToolCall{Name: "memory_forget", Arguments: map[string]interface{}{"id": "id"}}
resp, err := a.handleToolCall(ctx, "s1", tc)
assert.NoError(t, err)
assert.Equal(t, "Memory forgotten", resp)
})

t.Run("unknown tool", func(t *testing.T) {
a := NewAgent(nil, nil, nil, nil, nil, false)
tc := gemini.ToolCall{Name: "unknown"}
_, err := a.handleToolCall(ctx, "s1", tc)
assert.Error(t, err)
assert.Contains(t, err.Error(), "unknown tool")
})
}

func TestAgent_ConfirmAction(t *testing.T) {
t.Run("Non-Interactive", func(t *testing.T) {
a := &Agent{InteractiveMode: false}
assert.True(t, a.confirmAction("test"))
})

t.Run("Interactive Yes", func(t *testing.T) {
// Simulate 'y' input
oldStdin := os.Stdin
r, w, _ := os.Pipe()
os.Stdin = r
fmt.Fprintln(w, "y")
w.Close()
defer func() { os.Stdin = oldStdin }()

a := &Agent{InteractiveMode: true}
assert.True(t, a.confirmAction("test"))
})

t.Run("Interactive No", func(t *testing.T) {
// Simulate 'n' input
oldStdin := os.Stdin
r, w, _ := os.Pipe()
os.Stdin = r
fmt.Fprintln(w, "n")
w.Close()
defer func() { os.Stdin = oldStdin }()

a := &Agent{InteractiveMode: true}
assert.False(t, a.confirmAction("test"))
})
}
