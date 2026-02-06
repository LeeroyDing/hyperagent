package agent

import (
"context"
"fmt"
"log/slog"
"strings"

"github.com/LeeroyDing/hyperagent/internal/editor"
"github.com/LeeroyDing/hyperagent/internal/executor"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/mcp"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/LeeroyDing/hyperagent/internal/orchestrator"
"github.com/LeeroyDing/hyperagent/internal/token"
)

// GeminiClient defines the interface for interacting with Gemini
type GeminiClient interface {
GenerateContent(ctx context.Context, messages []gemini.Message) (string, error)
}

type Agent struct {
InteractiveMode bool
DryRun          bool
Gemini          GeminiClient
Executor        *executor.ShellExecutor
Memory          *memory.Memory
MCP             *mcp.MCPManager
History         *history.HistoryManager
TokenMgr        *token.TokenManager
Editor          *editor.FileEditor
Orchestrator    *orchestrator.Orchestrator
}

func NewAgent(gemini GeminiClient, executor *executor.ShellExecutor, memory *memory.Memory, mcpMgr *mcp.MCPManager, historyMgr *history.HistoryManager, interactiveMode bool) *Agent {
return &Agent{
Gemini:          gemini,
Executor:        executor,
Memory:          memory,
MCP:             mcpMgr,
History:         historyMgr,
InteractiveMode: interactiveMode,
Editor:          editor.NewFileEditor(),
Orchestrator:    orchestrator.NewOrchestrator(),
}
}

func (a *Agent) Run(ctx context.Context, sessionID, prompt string) (string, error) {
slog.Info("Starting agentic loop", "session", sessionID, "prompt", prompt)

hist, err := a.History.LoadHistory(sessionID)
if err != nil {
return "", fmt.Errorf("failed to load history: %w", err)
}

var messages []gemini.Message
for _, m := range hist {
messages = append(messages, gemini.Message{Role: m.Role, Content: m.Content})
}
messages = append(messages, gemini.Message{Role: "user", Content: prompt})

resp, err := a.Gemini.GenerateContent(ctx, messages)
if err != nil {
return "", fmt.Errorf("gemini error: %w", err)
}

return resp, nil
}

func (a *Agent) confirmAction(action string) bool {
if !a.InteractiveMode {
return true
}
fmt.Printf("\n[INTERACTIVE MODE] Confirm action: %s (y/n): ", action)
var response string
fmt.Scanln(&response)
return strings.ToLower(strings.TrimSpace(response)) == "y"
}
