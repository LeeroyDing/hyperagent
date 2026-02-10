package executor

import (
"fmt"
"testing"

"github.com/stretchr/testify/assert"
)

// MockShell implements the Shell interface for testing
type MockShell struct {
Responses map[string]string
Closed    bool
}

func (m *MockShell) Execute(command string) (string, error) {
if m.Closed {
return "", fmt.Errorf("session closed")
}
if resp, ok := m.Responses[command]; ok {
return resp, nil
}
return "mock output for " + command, nil
}

func (m *MockShell) Close() error {
m.Closed = true
return nil
}

func TestShellExecutor_WithMock(t *testing.T) {
t.Run("Allowlist and Mock Execution", func(t *testing.T) {
e := NewShellExecutor([]string{"ls", "echo"})

// Inject mock creator
mockShell := &MockShell{
Responses: map[string]string{
"echo hello": "hello",
},
}
e.Manager.Creator = func(id string) (Shell, error) {
return mockShell, nil
}

// Test allowed command
got, err := e.Execute("sess1", "echo hello")
assert.NoError(t, err)
assert.Equal(t, "hello", got)

// Test blocked command
_, err = e.Execute("sess1", "pwd")
assert.Error(t, err)
assert.Contains(t, err.Error(), "not in the allowlist")

// Test empty command
_, err = e.Execute("sess1", "")
assert.Error(t, err)
assert.Equal(t, "empty command", err.Error())
})

t.Run("Session Persistence", func(t *testing.T) {
e := NewShellExecutor(nil) // No allowlist

count := 0
e.Manager.Creator = func(id string) (Shell, error) {
count++
return &MockShell{}, nil
}

_, _ = e.Execute("sess1", "cmd1")
_, _ = e.Execute("sess1", "cmd2")
_, _ = e.Execute("sess2", "cmd1")

assert.Equal(t, 2, count, "Should have created exactly 2 sessions")
})

t.Run("Cleanup", func(t *testing.T) {
e := NewShellExecutor(nil)
mock := &MockShell{}
e.Manager.Creator = func(id string) (Shell, error) {
return mock, nil
}

_, _ = e.Execute("sess1", "cmd")
e.Cleanup()
assert.True(t, mock.Closed)
})

t.Run("Creator Error", func(t *testing.T) {
e := NewShellExecutor(nil)
e.Manager.Creator = func(id string) (Shell, error) {
return nil, fmt.Errorf("spawn failed")
}

_, err := e.Execute("sess1", "cmd")
assert.Error(t, err)
assert.Contains(t, err.Error(), "spawn failed")
})
}
