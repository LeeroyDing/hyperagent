package integration

import (
"context"
"testing"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/gemini"
)

func TestE2E_AgentLoop(t *testing.T) {
// 1. Setup Mocks
mockGemini := &MockGeminiClient{
Responses: []string{
"", // First response triggers a tool call
"The output is: Mock output for: echo 'hello'", // Second response after tool execution
},
ToolCalls: [][]gemini.ToolCall{
{
{Name: "execute_command", Arguments: map[string]interface{}{"command": "echo 'hello'"}},
},
nil,
},
}

mockExecutor := &MockExecutor{}
mockMemory := &MockMemory{}
mockHistory := &MockHistory{}

// 2. Initialize Agent
agent := agent.NewAgent(mockGemini, mockExecutor, mockMemory, nil, mockHistory, false)

// 3. Run Agent Loop
ctx := context.Background()
sessionID := "test-session"
prompt := "Say hello"

response, err := agent.Run(ctx, sessionID, prompt)
if err != nil {
t.Fatalf("Agent.Run failed: %v", err)
}

// 4. Verify Results
// Check if the tool was executed
if len(mockExecutor.ExecutedCommands) != 1 {
t.Errorf("Expected 1 command execution, got %d", len(mockExecutor.ExecutedCommands))
}
if len(mockExecutor.ExecutedCommands) > 0 && mockExecutor.ExecutedCommands[0] != "echo 'hello'" {
t.Errorf("Expected command 'echo 'hello'', got '%s'", mockExecutor.ExecutedCommands[0])
}

// Check the final response
expectedResponse := "The output is: Mock output for: echo 'hello'"
if response != expectedResponse {
t.Errorf("Expected response '%s', got '%s'", expectedResponse, response)
}

// Check history
history, _ := mockHistory.LoadHistory(sessionID)
if len(history) == 0 {
t.Error("Expected history to be recorded")
}
}
