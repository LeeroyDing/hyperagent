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
"github.com/google/generative-ai-go/genai"
)

// GeminiClient defines the interface for interacting with Gemini
type GeminiClient interface {
GenerateContent(ctx context.Context, messages []gemini.Message, tools []*genai.Tool) (string, []gemini.ToolCall, error)
SendToolResponse(ctx context.Context, messages []gemini.Message, tools []*genai.Tool, toolResponses []gemini.ToolResponse) (string, []gemini.ToolCall, error)
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

func (a *Agent) getTools() []*genai.Tool {
return []*genai.Tool{
{
FunctionDeclarations: []*genai.FunctionDeclaration{
{
Name:        "execute_command",
Description: "Execute a shell command on the host system",
Parameters: &genai.Schema{
Type: genai.TypeObject,
Properties: map[string]*genai.Schema{
"command": {Type: genai.TypeString, Description: "The shell command to execute"},
},
Required: []string{"command"},
},
},
{
Name:        "read_file",
Description: "Read lines from a file",
Parameters: &genai.Schema{
Type: genai.TypeObject,
Properties: map[string]*genai.Schema{
"path":  {Type: genai.TypeString, Description: "Path to the file"},
"start": {Type: genai.TypeInteger, Description: "Start line (1-indexed)"},
"end":   {Type: genai.TypeInteger, Description: "End line (optional)"},
},
Required: []string{"path"},
},
},
{
Name:        "replace_text",
Description: "Replace text in a file",
Parameters: &genai.Schema{
Type: genai.TypeObject,
Properties: map[string]*genai.Schema{
"path":     {Type: genai.TypeString, Description: "Path to the file"},
"old_text": {Type: genai.TypeString, Description: "Text to find"},
"new_text": {Type: genai.TypeString, Description: "Replacement text"},
},
Required: []string{"path", "old_text", "new_text"},
},
},
{
Name:        "memorize",
Description: "Save information to long-term memory",
Parameters: &genai.Schema{
Type: genai.TypeObject,
Properties: map[string]*genai.Schema{
"id":      {Type: genai.TypeString, Description: "Unique ID for the memory"},
"content": {Type: genai.TypeString, Description: "Content to memorize"},
},
Required: []string{"id", "content"},
},
},
{
Name:        "recall",
Description: "Search long-term memory",
Parameters: &genai.Schema{
Type: genai.TypeObject,
Properties: map[string]*genai.Schema{
"query": {Type: genai.TypeString, Description: "Search query"},
"limit": {Type: genai.TypeInteger, Description: "Max results (default 5)"},
},
Required: []string{"query"},
},
},
},
},
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

tools := a.getTools()
textResp, toolCalls, err := a.Gemini.GenerateContent(ctx, messages, tools)
if err != nil {
return "", fmt.Errorf("gemini error: %w", err)
}

for len(toolCalls) > 0 {
var toolResponses []gemini.ToolResponse
for _, tc := range toolCalls {
slog.Info("Handling tool call", "name", tc.Name, "args", tc.Arguments)
result, err := a.handleToolCall(ctx, tc)
if err != nil {
result = fmt.Sprintf("Error: %v", err)
}
toolResponses = append(toolResponses, gemini.ToolResponse{
Name:    tc.Name,
Content: result,
})
}

// Update messages with the assistant's tool calls and the responses
// Note: In a real implementation, we'd need to handle the history more carefully
// for the multi-turn tool calling loop.
textResp, toolCalls, err = a.Gemini.SendToolResponse(ctx, messages, tools, toolResponses)
if err != nil {
return "", fmt.Errorf("gemini tool response error: %w", err)
}
}

return textResp, nil
}

func (a *Agent) handleToolCall(ctx context.Context, tc gemini.ToolCall) (string, error) {
switch tc.Name {
case "execute_command":
cmd := tc.Arguments["command"].(string)
if !a.confirmAction(fmt.Sprintf("Execute command: %s", cmd)) {
return "Action cancelled by user", nil
}
return a.Executor.Execute(cmd)
case "read_file":
path := tc.Arguments["path"].(string)
start := int(tc.Arguments["start"].(float64))
end := 0
if e, ok := tc.Arguments["end"]; ok {
end = int(e.(float64))
}
lines, err := a.Editor.ReadLines(path, start, end)
if err != nil {
return "", err
}
return strings.Join(lines, "\n"), nil
case "replace_text":
path := tc.Arguments["path"].(string)
old := tc.Arguments["old_text"].(string)
new := tc.Arguments["new_text"].(string)
if !a.confirmAction(fmt.Sprintf("Replace text in %s", path)) {
return "Action cancelled by user", nil
}
err := a.Editor.Replace(path, old, new)
if err != nil {
return "", err
}
return "Text replaced successfully", nil
case "memorize":
id := tc.Arguments["id"].(string)
content := tc.Arguments["content"].(string)
err := a.Memory.Memorize(ctx, id, content, nil)
if err != nil {
return "", err
}
return "Information memorized", nil
case "recall":
query := tc.Arguments["query"].(string)
limit := 5
if l, ok := tc.Arguments["limit"]; ok {
limit = int(l.(float64))
}
results, err := a.Memory.Recall(ctx, query, limit)
if err != nil {
return "", err
}
var sb strings.Builder
for _, r := range results {
sb.WriteString(fmt.Sprintf("ID: %s\nContent: %s\n\n", r.ID, r.Content))
}
return sb.String(), nil
default:
return "", fmt.Errorf("unknown tool: %s", tc.Name)
}
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
