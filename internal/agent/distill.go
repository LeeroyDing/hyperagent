package agent

import (
"context"
"fmt"
"log/slog"

"github.com/LeeroyDing/hyperagent/internal/gemini"
)

func (a *Agent) Distill(ctx context.Context, sessionID string) error {
slog.Info("Starting memory distillation", "session", sessionID)

hist, err := a.History.LoadHistory(sessionID)
if err != nil {
return fmt.Errorf("failed to load history for distillation: %w", err)
}

if len(hist) < 5 {
return nil // Not enough context to distill
}

// Prepare history for Gemini
var historyText string
for _, m := range hist {
historyText += fmt.Sprintf("%s: %s\n", m.Role, m.Content)
}

prompt := fmt.Sprintf("Summarize the following conversation into a concise set of key facts, decisions, and context for long-term memory. Focus on information that will be useful for future interactions. Conversation:\n%s", historyText)

summary, _, err := a.Gemini.GenerateContent(ctx, []gemini.Message{
{Role: "user", Content: prompt},
}, nil)
if err != nil {
return fmt.Errorf("failed to generate distillation summary: %w", err)
}

// Save to vector memory
err = a.Memory.Memorize(ctx, fmt.Sprintf("distill-%s-%d", sessionID, len(hist)), summary, map[string]string{
"session_id": sessionID,
"type":       "distillation",
})
if err != nil {
return fmt.Errorf("failed to save distillation to memory: %w", err)
}

slog.Info("Memory distillation complete", "session", sessionID)
return nil
}
